package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Repository repository
type Repository interface {
	Add(key string, content string)
	List(key string) ([]string, error)

	GetHTTPContent(key string) (HttpContent, error)
	SetHTTPContent(key string, content HttpContent) error
}

// redisRepository repository using redis as backend
type redisRepository struct {
	prefixKey string
	ttl       int
	pool      *redis.Pool
}

// NewRedisRepository redisRepository constructor
func NewRedisRepository(addr string, ttl int) Repository {
	return &redisRepository{
		prefixKey: "rrp::",
		ttl:       ttl,
		pool:      createPool(addr),
	}
}

func (r *redisRepository) GetHTTPContent(key string) (HttpContent, error) {
	result, err := get(r.pool, r.prefixKey+key)
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
	set(r.pool, r.prefixKey+key, string(json), r.ttl)
	return nil
}

func (r *redisRepository) Add(key string, content string) {
	rpush(r.pool, r.prefixKey+key, content, r.ttl)
}

func (r *redisRepository) List(key string) ([]string, error) {
	return lrange(r.pool, r.prefixKey+key)
}

func lrange(redisPool *redis.Pool, key string) ([]string, error) {
	c := redisPool.Get()
	defer c.Close()

	values, err := redis.Strings(c.Do("LRANGE", key, 0, -1))
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return nil, errors.New("empty key")
	}
	return values, err
}

func get(redisPool *redis.Pool, key string) (string, error) {
	c := redisPool.Get()
	defer c.Close()

	value, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}
	if len(value) == 0 {
		return "", errors.New("empty key")
	}
	return value, err
}

func rpush(redisPool *redis.Pool, key, value string, ttl int) error {
	c := redisPool.Get()
	defer c.Close()

	c.Send("MULTI")
	c.Send("RPUSH", key, value)
	c.Send("EXPIRE", key, ttl)
	_, err := c.Do("EXEC")
	return err
}

func set(redisPool *redis.Pool, key, value string, ttl int) error {
	c := redisPool.Get()
	defer c.Close()

	_, err := c.Do("SETEX", key, ttl, value)
	return err
}

func createPool(addr string) *redis.Pool {
	Info.Println("initializing redis...")

	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
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
