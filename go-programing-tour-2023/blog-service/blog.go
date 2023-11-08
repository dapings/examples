package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/routers"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/logger"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/setting"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/tracer"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

// @title 博客系统
// @version 1.0
// @description Go 语言 blog service
// @termsOfService https://github.com/dapings/examples
func main() {
	global.Logger.Infof(context.TODO(), "%s : go-programming/blog", "blog-service")
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
	err := setupFlag()
	if err != nil {
		log.Fatalf("init.setupFlag err: %v", err)
	}
	err = setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	err = setupTracer()
	if err != nil {
		log.Fatalf("init.setupTracer err: %v", err)
	}
}

var (
	port      string
	runMode   string
	config    string
	isVersion bool
)

func setupSetting() error {
	// 也可以将配置文件存放在系统自带的全局变量中，如$HOME/conf或/etc/conf中，好处是不需要重新自定义一个新的系统环境变更
	// 内置一些系统环境变量的读取，优先级低于命令行参数，但高于文件配置。
	// 或者将配置文件打包到二进制文件中，通过 go-bindata 库可以将数据文件转换为Go代码，就可以摆脱静态资源文件了。
	// 或直接使用集中式的配置中心。
	s, err := setting.NewSetting(strings.Split(config, ",")...)
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
	err = s.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}
	err = s.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	global.AppSetting.DefaultContextTimeout *= time.Second
	global.JWTSetting.Expire *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second
	// 如果存在，则覆盖原有的文件配置
	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	return nil
}

func setupDBEngine() error {
	// 注：:= 会重新声明并创建左侧的新局部变量，因此在其他包调用变量时，仍然是 nil，达不到可用标准
	// 在赋值时并没有赋值到真正需要赋值的变量上
	var err error
	// global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupLogger() error {
	fileName := filepath.Join(global.AppSetting.LogSavePath, global.AppSetting.LogFileName+global.AppSetting.LogFileExt)
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500,
		MaxAge:     10,
		MaxBackups: 15,
		LocalTime:  true,
		Compress:   false,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func setupTracer() error {
	jaegerTracer, _, err := tracer.NewJaegerTracer("blog-service", "127.0.0.1:6831")
	if err != nil {
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func setupFlag() error {
	// 如果存在命令行参数，则优先使用命令行参数，否则使用配置文件中的配置参数
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径，以逗号(,)分隔")
	flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}
