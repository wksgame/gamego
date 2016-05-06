package base

import (
	"log"
	"os"
	"os/signal"
	//	"time"
)

var exitSignal chan os.Signal
var exitRegiste bool = false

//type ExitCallback func()

// 退出程序前调用cb
// type ExitCallback func()
func ExitApplication(cb func(interface{}), arg interface{}) bool {

	if exitRegiste {
		return false
	}

	exitRegiste = true

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

		//		for i := 3; i > 0; i-- {
		//			log.Printf("退出倒计时%d秒", i)
		//			time.Sleep(time.Second)
		//		}
	}()

	return true
}
