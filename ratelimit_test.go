package ratelimit

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMiddlewareFunc(t *testing.T) {
	r := gin.New()

	rules := []*Rule{
		// global rule
		&Rule{
			Global: true,
			Limit:  1,
		},
		// api rules
		&Rule{
			Method: "GET",
			Path:   "/hello",
			Limit:  1,
		},
		&Rule{
			Method: "GET",
			Path:   "/user/:id",
			Limit:  1,
		},
	}

	mw, err := New(rules...)
	if err != nil {
		panic(err)
	}

	r.Use(mw.MiddlewareFunc())

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, "hello")
	})
	r.GET("/user/:id", func(c *gin.Context) {
		c.JSON(200, c.Param("id"))
	})

	r.Run()
}
