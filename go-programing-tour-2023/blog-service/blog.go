package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/model"
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

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
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

func setupDBEngine() error {
	// 注：:= 会重新声明并创建左侧的新局部变量，因此在其他包调用变量时，仍然是 nil，达不到可用标准
	// 在赋值时并没有赋值到真正需要赋值的变量上
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}
