package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/routers"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/setting"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}

var (
	port      string
	runMode   string
	config    string
	isVersion bool
)

func setupSetting() error {
	s, err := setting.NewSetting("configs")
	if err != nil {
		return err
	}

	err = s.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	return nil
}
