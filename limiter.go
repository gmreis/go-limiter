package limiter

import (
	"time"

	"github.com/gin-gonic/gin"
)

// KeyDefine is a function used to defined the key of cache.
type KeyDefine func(*gin.Context) string

// AbortCallback is a function used to abort the request.
type AbortCallback func(*gin.Context)

// ConnectionCache define a methods used in the cache connection.
type ConnectionCache interface {
	Finish()
	Exist(string) (bool, error)
	Set(string, time.Duration) error
}

// CacheDriver define a method that return a cache connection.
type CacheDriver interface {
	GetConnection() ConnectionCache
}

// NewLimiter return a middleware used to limit a request.
func NewLimiter(cache CacheDriver, key KeyDefine, abort AbortCallback, limiter time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := cache.GetConnection()
		defer conn.Finish()

		k := key(c)
		ok, err := conn.Exist(k)
		if err != nil {
			panic(err)
		}

		if ok == false {
			c.Next()
			conn.Set(k, limiter)
		} else {
			abort(c)
		}
	}
}
