package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var currentTime = time.Now()

var DataDir string
var ConfigsDir string
var ServersDir string
var DownloadsDir string
var LogsDir string

func initDataDirs() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DataDir = filepath.Join(homePath, ".config", "MCSCS")
	year, month, day := currentTime.Date()
	ConfigsDir = filepath.Join(DataDir, "configs")
	ServersDir = filepath.Join(DataDir, "servers")
	DownloadsDir = filepath.Join(DataDir, "downloads")
	LogsDir = filepath.Join(DataDir, "logs", fmt.Sprintf("%d%02d%02d%02d", year, month, day, currentTime.Hour()))
	createDirIfNotExist(DataDir)
	createDirIfNotExist(ConfigsDir)
	createDirIfNotExist(ServersDir)
	createDirIfNotExist(DownloadsDir)
	createDirIfNotExist(LogsDir)
}

func createDirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
}
