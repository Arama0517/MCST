package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var currentTime = time.Now()

func initDataDirs() {
	Logger.Info("数据根目录: ", GetDataDir())
	Logger.Info("配置存放目录: ", GetConfigsDir())
	Logger.Info("服务器存放目录: ", GetServersDir())
	Logger.Info("Aria2c配置/可执行程序存放目录: ", GetAria2cDir())
	Logger.Info("日志存放目录: ", GetLogsDir())
}

// 获取存放数据的目录
func GetDataDir() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homePath, ".config", "MCSCS")
}

func GetConfigsDir() string {
	path := filepath.Join(GetDataDir(), "configs")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	return path
}

func GetServersDir() string {
	path := filepath.Join(GetDataDir(), "servers")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	return path
}

func GetAria2cDir() string {
	path := filepath.Join(GetDataDir(), "aria2c")
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	return path
}

func GetLogsDir() string {
	year, month, day := currentTime.Date()
	path := filepath.Join(GetDataDir(), "logs", fmt.Sprintf("%d%02d%02d%02d", year, month, day, currentTime.Hour()))
	err := os.MkdirAll(path, 0755)
	if err != nil {
		panic(err)
	}
	return path
}
