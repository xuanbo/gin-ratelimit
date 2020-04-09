package ratelimit

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// GinMiddleware token bucket implementation.
type GinMiddleware struct {
	globalBucket *ratelimit.Bucket
	buckets      map[string]*ratelimit.Bucket
}

// Rule the rule for rate limit.
type Rule struct {
	key    string
	Global bool   `json:"global" yaml:"global"`
	Method string `json:"method" yaml:"method"`
	Path   string `json:"path" yaml:"path"`
	Limit  int64  `json:"limit" yaml:"limit"`
}

// New for Middleware.
func New(rules ...*Rule) (*GinMiddleware, error) {
	mw := &GinMiddleware{
		globalBucket: nil,
		buckets:      make(map[string]*ratelimit.Bucket),
	}

	for _, rule := range rules {
		if rule.Limit <= 0 {
			return nil, fmt.Errorf("rule limit must great then 0")
		}

		if rule.Global {
			gb := mw.globalBucket
			if gb == nil {
				mw.globalBucket = ratelimit.NewBucket(time.Second, rule.Limit)
				continue
			} else {
				return nil, fmt.Errorf("global rule exist")
			}
		}

		if rule.Method == "" || rule.Path == "" {
			return nil, fmt.Errorf("method, path required")
		}
		key := rule.Method + "-" + rule.Path
		mw.buckets[key] = ratelimit.NewBucket(time.Second, rule.Limit)
	}
	return mw, nil
}

// MiddlewareFunc makes RateLimitMiddleware implement the Middleware interface.
func (mw GinMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ticket := mw.globalBucket.TakeAvailable(1); ticket == 0 {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		key := c.Request.Method + "-" + c.FullPath()
		bucket, ok := mw.buckets[key]
		if ok {
			if ticket := bucket.TakeAvailable(1); ticket == 0 {
				c.AbortWithStatus(http.StatusTooManyRequests)
				return
			}
		}

		c.Next()
	}
}
