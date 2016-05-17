package base

import (
	"log"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type RedisPool struct {
	pool    *pool.Pool
	dbindex int
}

func (self *RedisPool) dail(network, addr string) (*redis.Client, error) {
	redis_cli, err := redis.Dial(network, addr)
	if err != nil {
		log.Printf("RedisPool Dail error:%s", err)
		return nil, err
	}
	redis_cli.Cmd("select", self.dbindex)
	return redis_cli, err
}

func (self *RedisPool) Connect(host, port string, index, num int) bool {
	self.dbindex = index
	p, err := pool.NewCustom("tcp", host+":"+port, num, self.dail)

	if err != nil {
		log.Printf("RedisPool Connect error:%s", err)
		return false
	}

	self.pool = p

	log.Printf("redis connect ok,ip:%s port:%s dbindex:%d number:%d", host, port, index, num)
	return true
}

func (self *RedisPool) Get() *redis.Client {
	p, _ := self.pool.Get()
	return p
}

func (self *RedisPool) Put(c *redis.Client) {
	self.pool.Put(c)
}
