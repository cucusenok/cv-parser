package app

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func PGFullMatchQueryHandler(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, res)

}

func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/api/server/status", func(c *gin.Context) {
		retryCount := c.GetHeader("Retry-Count")
		userAgent := c.GetHeader("User-Agent")
		fmt.Println(userAgent)
		//c.JSON(http.StatusOK, gin.H{"serverIsOpen": true, "botname": "onchaincoin_bot", "code": 200})
		//return
		if retryCount == "9" {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		} else {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "pong",
			})
		}
	})

	r.GET("/", PGFullMatchQueryHandler)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
