package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmreis/go-limiter"
	"github.com/gmreis/go-limiter/drivers"
)

func getHashBody(c *gin.Context) string {
	body, _ := c.GetRawData()
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[0:])
}

func abortRequest(c *gin.Context) {
	c.AbortWithStatus(429)
}

func main() {
	r := gin.Default()
	redis := drivers.NewRedis("localhost:6379")

	limiter := limiter.NewLimiter(redis, getHashBody, abortRequest, 10*time.Second)

	r.POST("/v1/products",
		limiter,
		func(c *gin.Context) {
			c.Status(http.StatusNoContent)
		})

	r.Run(":8080")
}
