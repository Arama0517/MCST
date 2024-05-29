/*
 * MCSCS can be used to easily create, launch, and configure a Minecraft server.
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
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

var DataDir string
var serversDir string
var downloadsDir string
var logsDir string

var configsPath string

func InitData() {
	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DataDir = filepath.Join(UserHomeDir, ".config", "MCSCS")
	serversDir = filepath.Join(DataDir, "servers")
	downloadsDir = filepath.Join(DataDir, "downloads")
	logsDir = filepath.Join(DataDir, "logs")
	createDirIfNotExist(DataDir)
	createDirIfNotExist(serversDir)
	createDirIfNotExist(downloadsDir)
	createDirIfNotExist(logsDir)
	configsPath = filepath.Join(DataDir, "configs.json")
	if _, err := os.Stat(configsPath); os.IsNotExist(err) {
		jsonData, err := json.MarshalIndent(MCSCSConfig{
			LogLevel:      "info",
			API:           0,
			Downloads:     []DownloadInfo{},
			Javas:         []JavaInfo{},
			Servers:       map[string]ServerConfig{},
			MaxConnetions: 8,
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

type Ram struct {
	// Java 虚拟机初始堆内存
	XMX uint64 `json:"xmx"`
	// Java 虚拟机最大堆内存
	XMS uint64 `json:"xms"`
}

type ServerInfo struct {
	ServerType       string `json:"server_type"`
	MinecraftVersion string `json:"mc_version"`
	BuildVersion     string `json:"build_version"`
}

type ServerConfig struct {
	// 服务器名称
	Name string `json:"name"`

	// 服务器内存配置
	Ram Ram `json:"ram"`

	// 编码格式
	Encoding string `json:"encoding"`

	// Java
	Java JavaInfo `json:"java"`

	// 服务器核心信息
	Info ServerInfo `json:"info"`

	// Java虚拟机其他参数
	JVMArgs []string `json:"jvm_args"`

	// Minecraft服务器参数
	ServerArgs []string `json:"server_args"`
}

type DownloadInfo struct {
	Info ServerInfo `json:"info"`
	Path string     `json:"path"`
}

type MCSCSConfig struct {
	LogLevel      string                  `json:"log_level"` // 日志级别
	API           int                     `json:"api"`       // 使用的API, 0: 无极镜像, 1: 极星云镜像
	Downloads     []DownloadInfo          `json:"downloads"` // 下载列表
	Javas         []JavaInfo              `json:"javas"`     // Java列表
	Servers       map[string]ServerConfig `json:"servers"`   // 服务器列表, 如果服务器名称(key)为temp, CreatePage调用时会视为暂存配置而不是名为temp的服务器
	MaxConnetions int                     `json:"max_connections"` // 使用 Downloader{}.Download() 多线程下载时的最大连接数
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
	log.Info().Interface("config", config).Msg("加载配置")
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
	log.Info().Interface("config", c).Msg("保存配置")
	return nil
}
