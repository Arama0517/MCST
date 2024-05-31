/*
 * Minecraft Server Tool(MCST) is a command-line utility making Minecraft server creation quick and easy for beginners.
 * Copyright (C) 2024 Arama
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package lib

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
)

var DataDir string
var ServersDir string
var DownloadsDir string

var configsPath string

func InitData() {
	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DataDir = filepath.Join(UserHomeDir, ".config", "MCST")
	ServersDir = filepath.Join(DataDir, "servers")
	DownloadsDir = filepath.Join(DataDir, "downloads")
	createDirIfNotExist(DataDir)
	createDirIfNotExist(ServersDir)
	createDirIfNotExist(DownloadsDir)
	configsPath = filepath.Join(DataDir, "configs.json")
	if _, err := os.Stat(configsPath); os.IsNotExist(err) {
		jsonData, err := json.MarshalIndent(MCSCSConfig{
			LogLevel:       "info",
			API:            0,
			Cores:          []Core{},
			Servers:        map[string]ServerConfig{},
			MaxConnections: 8,
		}, "", "  ")
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(configsPath, jsonData, 0644)
		if err != nil {
			panic(err)
		}
	}
}

type Core struct {
	URL        url.URL     `json:"url"`
	FileName   string      `json:"file_name"`
	FilePath   string      `json:"file_path"`
	ExtrasData interface{} `json:"extras_data"`
}

type Java struct {
	Path     string   `json:"path"`
	Args     []string `json:"args"`
	Xmx      uint64   `json:"Xmx"`
	Xms      uint64   `json:"Xms"`
	Encoding string   `json:"encoding"`
}
type ServerConfig struct {
	// 服务器名称
	Name string `json:"name"`

	// Java
	Java Java `json:"java"`

	// Minecraft服务器参数
	ServerArgs []string `json:"server_args"`
}

type MCSCSConfig struct {
	LogLevel       string                  `json:"log_level"`       // 日志级别
	API            int                     `json:"api"`             // 使用的API, 0: 无极镜像, 1: 极星云镜像
	Cores          []Core                  `json:"cores"`           // 核心列表
	Servers        map[string]ServerConfig `json:"servers"`         // 服务器列表, 如果服务器名称(key)为temp, CreatePage调用时会视为暂存配置而不是名为temp的服务器
	MaxConnections int                     `json:"max_connections"` // 使用 Downloader{}.Download() 多线程下载时的最大连接数
}

func createDirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func LoadConfigs() (MCSCSConfig, error) {
	file, err := os.ReadFile(configsPath)
	if err != nil {
		return MCSCSConfig{}, err
	}
	var config MCSCSConfig
	err = json.Unmarshal(file, &config)
	if err != nil {
		return MCSCSConfig{}, err
	}
	return config, nil
}

func (c *MCSCSConfig) Save() error {
	jsonConfig, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configsPath, jsonConfig, 0644)
	if err != nil {
		return err
	}
	return nil
}
