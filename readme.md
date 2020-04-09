# gin-ratelimit

a ratelimit middleware for gin.

## Install

```shell
go get -v github.com/xuanbo/gin-ratelimit
```

## Usage

```go
package main

import (
	ratelimit "github.com/xuanbo/gin-ratelimit"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// 1. rate limit rules
	rules := []*ratelimit.Rule{
		// globa rule
		&ratelimit.Rule{
			Global: true,
			Limit:  1,
		},
		// api rules
		&ratelimit.Rule{
			Method: "GET",
			Path:   "/hello",
			Limit:  1,
		},
		&ratelimit.Rule{
			Method: "GET",
			Path:   "/user/:id",
			Limit:  1,
		},
	}

	// 2. ratelimit middleware
	mw, err := ratelimit.New(rules...)
	if err != nil {
		panic(err)
	}

	// 3. use middleware
	r.Use(mw.MiddlewareFunc())

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, "hello")
	})
	r.GET("/user/:id", func(c *gin.Context) {
		c.JSON(200, c.Param("id"))
	})

	r.Run()
}
```

## Thanks

- github.com/gin-gonic/gin
- github.com/juju/ratelimit
