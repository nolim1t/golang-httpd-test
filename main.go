package main

import (
    // System Libraries
    "fmt"

    // External libraries
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/gzip"
)


func info(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "pong",
    })
}

func main() {
    router := gin.Default()
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://nolim1t.co"},
        AllowHeaders:     []string{"Origin"},
        ExposeHeaders:    []string{"Content-Length"},
    }))
    router.Use(gzip.Gzip(gzip.DefaultCompression))

    r := router.Group("/api")
    r.GET("/info", info)

    err := router.Run(fmt.Sprintf(":%d", 3000))
    if err != nil {
        panic(err)
    }
}
