package lib

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

func Init() {
	// 数据
	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DataDir = filepath.Join(UserHomeDir, ".config", "MCSCS")
	ServersDir = filepath.Join(DataDir, "servers")
	DownloadsDir = filepath.Join(DataDir, "downloads")
	LogsDir = filepath.Join(DataDir, "logs")
	createDirIfNotExist(DataDir)
	createDirIfNotExist(ServersDir)
	createDirIfNotExist(DownloadsDir)
	createDirIfNotExist(LogsDir)
	LogFilePath = filepath.Join(LogsDir, time.Now().Format("2006010215")+".log")
	LogFile, err = os.Create(LogFilePath)
	if err != nil {
		panic(err)
	}
	// 如果不先删除之前创建的软链接会导致创建报错!!!
	err = os.Remove(filepath.Join(LogsDir, "latest.log"))
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	err = os.Symlink(LogFilePath, filepath.Join(LogsDir, "latest.log"))
	if err != nil {
		panic(err)
	}
	MCSCSConfigsPath = filepath.Join(DataDir, "configs.json")
	if _, err := os.Stat(MCSCSConfigsPath); os.IsNotExist(err) {
		jsonData, err := json.MarshalIndent(MCSCSConfig{
			LogLevel:    "info",
			API:         0,
			Downloads:   []DownloadInfo{},
			Javas:       []JavaInfo{},
			Servers:     map[string]ServerConfig{},
			Concurrency: 8,
		}, "", "  ")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(MCSCSConfigsPath, jsonData, 0644)
		if err != nil {
			panic(err)
		}
	}

	// 日志
	configs, err := LoadConfigs()
	if err != nil {
		panic(err)
	}
	LogLevel, err := logrus.ParseLevel(configs.LogLevel)
	if err != nil {
		LogLevel = logrus.InfoLevel
	}
	Logger = &logrus.Logger{
		Out: LogFile,
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
		ReportCaller: true,
		Level:        LogLevel,
	}

	Logger.Info("Hello world!")
	Logger.Info("本程序遵循GPLv3协议开源")
	Logger.Info("作者: Arama 3584075812@qq.com")
}
