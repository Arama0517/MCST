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
var aria2cDir string

var configsPath string

func InitData() {
	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DataDir = filepath.Join(UserHomeDir, ".config", "MCST")
	ServersDir = filepath.Join(DataDir, "servers")
	DownloadsDir = filepath.Join(DataDir, "downloads")
	aria2cDir = filepath.Join(DataDir, "aria2c")
	createDirIfNotExist(DataDir)
	createDirIfNotExist(ServersDir)
	createDirIfNotExist(DownloadsDir)
	createDirIfNotExist(aria2cDir)
	configsPath = filepath.Join(DataDir, "configs.json")
	if _, err := os.Stat(configsPath); os.IsNotExist(err) {
		jsonData, err := json.MarshalIndent(MCSCSConfig{
			Cores:          []Core{},
			Servers:        map[string]Server{},
			Aria2c:         Aria2c{
				RetryWait:        2,
				Split:            10,
				MaxConnectionPerServer: 16,
				MinSplitSize:     "5M",
			},
			}, "", "    ")
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
	URL        url.URL     `json:"url"`         // 下载地址(如果不是本地的话)
	FileName   string      `json:"file_name"`   // 文件名
	FilePath   string      `json:"file_path"`   // 文件路径
	ExtrasData interface{} `json:"extras_data"` // 其他数据
}

type Java struct {
	Path     string   `json:"path"`     // Java路径
	Args     []string `json:"args"`     // Java虚拟机参数
	Xmx      uint64   `json:"Xmx"`      // Java虚拟机最大堆内存
	Xms      uint64   `json:"Xms"`      // Java虚拟机初始堆内存
	Encoding string   `json:"encoding"` // 编码
}
type Server struct {
	Name       string   `json:"name"`        // 服务器名称
	Java       Java     `json:"java"`        // Java
	ServerArgs []string `json:"server_args"` // Minecraft服务器参数
}

type Aria2c struct {
	RetryWait int `json:"retry_wait"` // 重试等待时间(秒)
	Split int `json:"split"` // 分块大小(M)
	MaxConnectionPerServer int `json:"max_connection_per_server"` // 单服务器最大连接数
	MinSplitSize string `json:"min_split_size"` // 最小分块大小
}

type MCSCSConfig struct {
	Cores          []Core            `json:"cores"`           // 核心列表
	Servers        map[string]Server `json:"servers"`         // 服务器列表, 如果服务器名称(key)为temp, CreatePage调用时会视为暂存配置而不是名为temp的服务器
	Aria2c         Aria2c            `json:"aria2c"`          // aria2c配置
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
