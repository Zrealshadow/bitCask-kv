package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	bc "github.com/bitCaskKV/bitCask"
	"github.com/bitCaskKV/global"
	"github.com/bitCaskKV/pkg/logger"
	"github.com/bitCaskKV/pkg/setting"
	"github.com/bitCaskKV/server/routers"
	"gopkg.in/natefinch/lumberjack.v2"
)

// A simple Router
func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("Setting Config Load Fail Error:%s", err.Error())
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("Global Logger Load Fail Error:%s", err.Error())

	}

	// if global.ServerSetting.RunMode != "debug" {
	err = setupBitCaskEngine()
	if err != nil {
		log.Fatalf("Default Engine Initialize Fail Error:%s", err.Error())
	}
	// }
}

func main() {
	// serverMux := http.NewServeMux()
	//Get
	// serverMux.HandleFunc("/v1/db", func(rw http.ResponseWriter, r *http.Request) {
	// 	// r.Form.Get("Method")
	// 	method := r.Header.Get("USE")
	// 	global.Logger.Infof("Get Method %s", method)
	// 	switch method {
	// 	case "Put":
	// 		_, _ = rw.Write([]byte(method))
	// 	case "Get":
	// 		_, _ = rw.Write([]byte(method))
	// 	case "Del":
	// 		_, _ = rw.Write([]byte(method))
	// 	default:
	// 		_, _ = rw.Write([]byte("HelloWorld"))
	// 	}
	// })
	// r := gin.Default()
	// r.GET("/v1/db", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "ok"})
	// })
	r := routers.NewRouter()
	r.Run(global.ServerSetting.Addr + ":" + global.ServerSetting.HttpPort)
}

func setupSetting() error {
	setting, _ := setting.NewSetting()

	err := setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return errors.New(fmt.Sprintf("Server Config Read Error: %s", err.Error()))
	}

	err = setting.ReadSection("Logger", &global.LoggerSetting)
	if err != nil {
		return errors.New(fmt.Sprintf("Logger Config Read Error: %s", err.Error()))
	}

	err = setting.ReadSection("DB", &global.BitCaskSetting)
	if err != nil {
		return errors.New(fmt.Sprintf("BitCask Engine Config Read Error: %s", err.Error()))
	}

	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(
		&lumberjack.Logger{
			Filename:  global.LoggerSetting.LogSavePath + string(os.PathSeparator) + global.LoggerSetting.LogFileExt,
			MaxSize:   global.LoggerSetting.MaxPageSize,
			MaxAge:    global.LoggerSetting.DefaultPageSize,
			LocalTime: true,
		}, "", log.LstdFlags)
	return nil
}

func setupBitCaskEngine() error {
	global.DefaultBitCaskEngine, _ = bc.NewBitCaskEngine(global.BitCaskSetting.MountDir)
	return nil
}
