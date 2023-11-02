package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSONP(http.StatusOK, gin.H{"message": "pong"})
	})
	_ = r.Run()
}

var (
	port      string
	runMode   string
	config    string
	isVersion bool
)
