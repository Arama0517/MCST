package lib

import (
	"os"
	"path/filepath"
)

// 获取存放数据的目录
func GetDataDir() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(homePath, ".config", "MCSCS")
}

func GetConfigsDir() string {
	return filepath.Join(GetDataDir(), "configs")
}

func GetServersDir() string {
	return filepath.Join(GetDataDir(), "servers")
}

func GetAria2cDir() string {
	return filepath.Join(GetDataDir(), "aria2c")
}

func GetLogsDir() string {
	return filepath.Join(GetDataDir(), "logs")
}
