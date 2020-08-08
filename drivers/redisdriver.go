package drivers

import (
	"time"

	"github.com/gmreis/go-limiter"
	"github.com/gomodule/redigo/redis"
)

// RedisDriver ...
type RedisDriver struct {
	pool *redis.Pool
}

// RedisConnection ...
type RedisConnection struct {
	conn redis.Conn
}

// NewRedis ....
func NewRedis(server string) limiter.CacheDriver {
	var pool = redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	// TODO Break em caso de Erro!!!

	return &RedisDriver{
		pool: &pool,
	}
}

// GetConnection ...
func (cache *RedisDriver) GetConnection() limiter.ConnectionCache {
	return &RedisConnection{
		conn: cache.pool.Get(),
	}
}

// Finish ...
func (ctx *RedisConnection) Finish() {
	ctx.conn.Close()
}

// Exist ...
func (ctx *RedisConnection) Exist(key string) (bool, error) {
	return redis.Bool(ctx.conn.Do("EXISTS", key))
}

// Set ...
func (ctx *RedisConnection) Set(key string, limiter time.Duration) error {
	_, err := redis.String(ctx.conn.Do("SET", key, true, "EX", limiter.Seconds()))
	return err
}
