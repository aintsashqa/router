# Router

### Install:

```cmd
go get -u github.com/aintsashqa/router
```

### Example:

```golang
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aintsashqa/router"
)

func handleHelloWorld(ctx *router.Context) error {
	return ctx.Plain(http.StatusOK, "Hello, world!")
	// return ctx.Json(http.StatusOK, map[string]string{"message": "Hello, world!"})
}

func handleHelloName(ctx *router.Context) error {
	name := ctx.Param("name")
	return ctx.Plain(http.StatusOK, fmt.Sprintf("Hello, %s!", name))
	// return ctx.Json(http.StatusOK, map[string]string{"message": fmt.Sprintf("Hello, %s!", name)})
}

func handleLogMiddleware(next router.HandlerFunc) router.HandlerFunc {
	return func(ctx *router.Context) error {
		log.Print("[", ctx.Request().Method, "] ", ctx.Request().URL.Path)
		return next(ctx)
	}
}

func main() {
	r := router.New()
	r.Use(handleLogMiddleware)
	r.Get("/", handleHelloWorld)
	r.Get("/:name", handleHelloName)

	http.ListenAndServe(":8080", &r)
}
```
