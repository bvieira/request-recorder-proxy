package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

var BASEKEY = "rrp::"

var redisPool = getPool()

func GetHttpContent(key string) (HttpContent, error) {
	if result, err := Get(key); err != nil {
		return HttpContent{}, err
	} else {
		var content HttpContent
		err2 := json.Unmarshal([]byte(result), &content)
		return content, err2
	}
}

func SetHttpContent(key string, content HttpContent) error {
	if cacheJson, err := json.Marshal(content); err != nil {
		return err
	} else {
		Set(key, string(cacheJson))
		return nil
	}
}

func Set(key string, content string) {
	setKey(BASEKEY+key, content)
}

func Get(key string) (string, error) {
	return getKey(BASEKEY + key)
}

func getKey(key string) (string, error) {
	c := redisPool.Get()
	defer c.Close()

	value, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}
	if len(value) == 0 {
		return "", errors.New("empty cache")
	}
	return value, err
}

func setKey(key string, value string) bool {
	c := redisPool.Get()
	defer c.Close()

	_, err := c.Do("SET", key, value)
	_, err2 := c.Do("EXPIRE", key, 2*60*60)
	return err != nil && err2 != nil
}

func getPool() *redis.Pool {
	Info.Println("initializing redis...")

	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				Error.Printf("err:[%v]\n", err)
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
