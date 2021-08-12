package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wolf88804/blog-service/global"
	"github.com/wolf88804/blog-service/internal/model"
	"github.com/wolf88804/blog-service/internal/routers"
	"github.com/wolf88804/blog-service/pkg/logger"
	"github.com/wolf88804/blog-service/pkg/setting"
	"github.com/wolf88804/blog-service/pkg/tracer"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	setupFlag()
	//初始化配置
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
	//初始化数据库
	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}
	//初始化日志
	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	//初始化追踪
	err=setupTracer()
	if err!=nil{
		log.Fatalf("init.setupTracer err: %v",err)
	}

}
var(
	port string
	runMode string
	config string
	isVersion bool

	buildTime string
	buildVersion string
	gitCommitID string
)
func setupFlag() {
	flag.StringVar(&port,"port","8000","启动端口")
	flag.StringVar(&runMode,"mode","debug","启动模式")
	flag.StringVar(&config,"config","configs/","指定要使用的配置文件路径")
	flag.BoolVar(&isVersion,"version",false,"编译信息")
	flag.Parse()
}

func setupTracer() error {
	jaegerTracer,_,err:=tracer.NewJaegerTracer("blog-service","127.0.0.1:6831")
	if err!=nil{
		return err
	}
	global.Tracer = jaegerTracer
	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)
	return nil
}

func setupDBEngine() error {
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupSetting() error {
	setting, err := setting.NewSetting(strings.Split(config,",")...)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	err = setting.ReadSection("Email", &global.EmailSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout += time.Second

	global.JWTSetting.Expire *= time.Second

	if port!=""{
		global.ServerSetting.HttpPort = port
	}

	if runMode!=""{
		global.ServerSetting.RunMode = runMode
	}

	return nil
}

// @title 博客系统
// @version 1.0
// @description Go 编程
// @termsOfService https://github.com
func main() {
	if isVersion{
		fmt.Printf("build_time: %s\n",buildTime)
		fmt.Printf("build_version: %s\n",buildVersion)
		fmt.Printf("git_commit_id: %s\n",gitCommitID)
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

	go func() {
		if err:= s.ListenAndServe();err!=nil && err!=http.ErrServerClosed{
			log.Fatalf("s.ListenAndServe err: %v",err)
		}
	}()

	quit:=make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit

	log.Println("Shuting down server...")

	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	if err:=s.Shutdown(ctx);err!=nil{
		log.Fatalf("Server forced to shutdown: %v",err)
	}
	log.Println("Server exiting")
}
