package base

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
)

type Once2 struct {
	m    sync.Mutex
	done uint32
}

func (o *Once2) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		panic("repeat")
		return
	}

	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	} else {
		panic("repeat2")
	}
}

var exitSignal chan os.Signal
var exitOnce Once2

// 退出程序前调用cb,仅能调用一次
func ExitApplication(cb func(interface{}), arg interface{}) {
	exitOnce.Do(func() {
		exitSignal = make(chan os.Signal, 1)
		signal.Notify(exitSignal, os.Interrupt, os.Kill)

		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
				}
				os.Exit(1)
			}()

			s := <-exitSignal

			log.Println("recv signal:", s)

			if cb != nil {
				cb(arg)
			}
		}()
	})
}
