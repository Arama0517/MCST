package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
	year, month, day := currentTime.Date()
	LogsDir = filepath.Join(DataDir, "logs", fmt.Sprintf("%d%02d%02d%02d", year, month, day, currentTime.Hour()))
	createDirIfNotExist(DataDir)
	createDirIfNotExist(ServersDir)
	createDirIfNotExist(DownloadsDir)
	createDirIfNotExist(LogsDir)
	MCSCSConfigsPath = filepath.Join(DataDir, "configs.json")
	if _, err := os.Stat(MCSCSConfigsPath); os.IsNotExist(err) {
		jsonData, err := json.MarshalIndent(MCSCSConfig{
			LogLevel:  "info",
			API:       0,
			Downloads: []DownloadInfo{},
			Javas:     []JavaInfo{},
			Servers:   map[string]ServerConfig{},
		}, "", "  ")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(MCSCSConfigsPath, jsonData, 0644)
		if err != nil {
			panic(err)
		}
	}

	// Logger
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{})
	logFile, err := os.OpenFile(filepath.Join(LogsDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Logger.SetOutput(logFile)
	Logger.Info("Hello world!")
	Logger.Info("本程序遵循GPLv3协议开源")
	Logger.Info("作者: Arama 3584075812@qq.com")

}
