# Gomisc

## LOG

````
package main

import (
	"gomisc/log"
	"gomisc/tcp"
)

func main() {
	log.Setup("gomisc", "10.141.48.10:4560")
}

````

## TCP

````
package main

import (
	"gomisc/log"
	"gomisc/tcp"
)

func main() {
	tcp.ListenAndServeWithSignal(":8000", &tcp.EchoHandler{})
}

````
## MVC

````

package main

import (
	"github.com/leehoawki/gomvc"
	"net/http"
)

func main() {
	r := mvc.New()
	r.Use(mvc.Recovery())
	r.Use(mvc.Logger())

	r.GET("/", func(c *mvc.Context) {
		c.HTML(http.StatusOK, "<h1>Hello</h1>")
	})

	r.GET("/hello", func(c *mvc.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/panic", func(c *mvc.Context) {
		names := []string{"test"}
		c.String(http.StatusOK, names[100])
	})

	r.POST("/login", func(c *mvc.Context) {
		c.JSON(http.StatusOK, mvc.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":9999")
}