package main

import (
	"net/http"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/routers"
)

func main() {
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}

var (
	port      string
	runMode   string
	config    string
	isVersion bool
)
