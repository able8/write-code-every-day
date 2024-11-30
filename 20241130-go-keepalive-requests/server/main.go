package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

type ConnRequest struct {
	requests uint32
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation.
type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "net/http context value " + k.name }

func NewConnRequest() *ConnRequest {
	return &ConnRequest{}
}

func (c *ConnRequest) Requests() uint32 {
	return atomic.AddUint32(&c.requests, 1)
}

var connRequestContextKey = &contextKey{"connRequest"}

func KeepAliveRequests(requests uint32) gin.HandlerFunc {
	return func(c *gin.Context) {
		obj := c.Request.Context().Value(connRequestContextKey)
		if cr, ok := obj.(*ConnRequest); ok {
			connRequests := cr.Requests()
			if connRequests >= requests {
				c.Header("Connection", "close")
			}
			slog.InfoContext(c.Request.Context(), "connRequest", slog.String("remoteAddr", c.Request.RemoteAddr), slog.Int("requests", int(connRequests)))
		}
		c.Next()
	}
}

func main() {
	r := gin.New()
	r.Use(KeepAliveRequests(10))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "hello world"})
	})
	svr := http.Server{
		Addr:              ":8888",
		Handler:           r,
		ReadTimeout:       60 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       65 * time.Second,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, connRequestContextKey, NewConnRequest())
		},
	}
	err := svr.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56520 requests=9
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56520 requests=10
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=1
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=2
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=3
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=4
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=5
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=6
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=7
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=8
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=9
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56521 requests=10
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=1
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=2
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=3
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=4
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=5
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=6
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=7
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=8
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=9
// 2024/11/30 12:13:55 INFO connRequest remoteAddr=[::1]:56523 requests=10
