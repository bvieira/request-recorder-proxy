package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Repository repository
type Repository interface {
	Set(key string, content string)
	Get(key string) (string, error)

	GetHTTPContent(key string) (HttpContent, error)
	SetHTTPContent(key string, content HttpContent) error
}

// redisRepository repository using redis as backend
type redisRepository struct {
	prefixKey string
	pool      *redis.Pool
}

// NewRedisRepository redisRepository constructor
func NewRedisRepository() Repository {
	return &redisRepository{
		prefixKey: "rrp::",
		pool:      createPool(),
	}
}

func (r *redisRepository) GetHTTPContent(key string) (HttpContent, error) {
	result, err := r.Get(key)
	if err != nil {
		return HttpContent{}, err
	}
	var content HttpContent
	return content, json.Unmarshal([]byte(result), &content)
}

func (r *redisRepository) SetHTTPContent(key string, content HttpContent) error {
	json, err := json.Marshal(content)
	if err != nil {
		return err
	}
	r.Set(key, string(json))
	return nil
}

func (r *redisRepository) Set(key string, content string) {
	setKey(r.pool, r.prefixKey+key, content, 2*60*60)
}

func (r *redisRepository) Get(key string) (string, error) {
	return getKey(r.pool, r.prefixKey+key)
}

func getKey(redisPool *redis.Pool, key string) (string, error) {
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

func setKey(redisPool *redis.Pool, key, value string, ttl int) error {
	c := redisPool.Get()
	defer c.Close()

	_, err := c.Do("SETEX", key, ttl, value)
	return err
}

func createPool() *redis.Pool {
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
