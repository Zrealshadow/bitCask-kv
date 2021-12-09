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

// @title 基于BitCask的K-V储存系统
// @version 1.0
// @description 基于BitCask K-V储存系统服务端
// @termsofService https://github.com/Zrealshadow/bitCask-kv
func main() {
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
