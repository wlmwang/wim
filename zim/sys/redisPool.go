package sys

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

func NewRedisPool() *redis.Pool {
	server := BaseConf.Get("redis").Get("host").MustString() + ":" + BaseConf.Get("redis").Get("port").MustString()
	password := BaseConf.Get("redis").Get("password").MustString()
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

var RedisPool *redis.Pool
