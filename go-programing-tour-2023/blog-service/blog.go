package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/dapings/examples/go-programing-tour-2023/blog-service/global"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/internal/routers"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/logger"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/setting"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/tracer"
	"github.com/dapings/examples/go-programing-tour-2023/blog-service/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/natefinch/lumberjack.v2"
)

// @title 博客系统
// @version 1.0
// @description Go 语言 blog service
// @termsOfService https://github.com/dapings/examples
func main() {
	topCtx := context.Background()
	global.Logger.Infof(topCtx, "%s : go-programming/blog", "blog-service")
	if isVersion {
		fmt.Printf("build_time: %s\n", buildTime)
		fmt.Printf("build_version: %s\n", buildVersion)
		fmt.Printf("git_commit_id: %s\n", gitCommitID)
		return
	}
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// 通过信号量的方式来解决问题：优雅重启和停止
	// 信号量是一种异步通知机制，用来提醒进程一个事件(硬件异常、程序执行异常、外部发生信息)已发生。如果进程定义的信号的处理函数，那么它将被执行，否则执行默认的处理函数。
	// kill -l 查看系统所支持的所有信息，SIGINT 进程结束，SIGTSTP 进程挂起, SIGQUIT 进程结束和dump core, SIGKILL 进程中断, SIGHUP 重启, SIGTERM 停止接收新请求
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.Logger.Fatalf(topCtx, "http server ListenAndServe err: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 最大时间控制，用于通知服务端它有 5s 的时间来处理原有的请求
	ctx, cancel := context.WithTimeout(topCtx, global.AppSetting.DefaultContextTimeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		global.Logger.Fatalf(ctx, "http server forced to shutdown, err: %v", err)
	}
	global.Logger.Info(ctx, "http server exiting")
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

	err = setupValidator()
	if err != nil {
		log.Fatalf("init.setupValidator err: %v", err)
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
	// -ldflags -X 参数，将信息写入变量中，格式为 package_name.variable_name=value
	// 构建信息：CGO_ENABLED=0 GOOS=linux go build -a -o blog-service . -ldflags="-w -s" -gcflags="-m -l" "-X main.buildTime=`data +%Y-%m-%d %H%M%S` -X main.buildVersion=v1.0.1 -X main.gitCommitID=`git rev-parse HEAD`"
	buildTime    string
	buildVersion string
	gitCommitID  string
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

func setupValidator() error {
	global.Validator = validator.NewCustomValidator()
	global.Validator.Engine()
	binding.Validator = global.Validator
	return nil
}
