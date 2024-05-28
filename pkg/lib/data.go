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
)

var DataDir string
var ConfigsDir string
var ServersDir string
var DownloadsDir string
var LogsDir string

var MCSCSConfigsPath string

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

type MCSCSConfig struct {
	// 日志级别
	LogLevel string `json:"log_level"`

	// 使用的API, 0: 无极镜像, 1: 极星云镜像
	API int `json:"api"`

	// 下载列表, 用`SaveDownloadsLists`函数更改
	Downloads []DownloadInfo `json:"downloads"`

	// Java列表, 用`SaveJavaLists`函数更改
	Javas []JavaInfo `json:"javas"`

	// 服务器列表, 用`SaveServerConfigs`函数更改
	Servers map[string]ServerConfig `json:"servers"`

	Concurrency int `json:"concurrency"`
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
	file, err := os.ReadFile(MCSCSConfigsPath)
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

func SaveConfigs(config MCSCSConfig) error {
	jsonConfig, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(MCSCSConfigsPath, jsonConfig, 0644)
	if err != nil {
		return err
	}
	return nil
}

func LoadServerConfigs() (map[string]ServerConfig, error) {
	configs, err := LoadConfigs()
	if err != nil {
		return nil, err
	}
	return configs.Servers, nil
}

func SaveServerConfigs(configs map[string]ServerConfig) error {
	MCSCSConfigs, err := LoadConfigs()
	if err != nil {
		return err
	}
	MCSCSConfigs.Servers = configs
	err = SaveConfigs(MCSCSConfigs)
	if err != nil {
		return err
	}
	return nil
}

func LoadJavaLists() ([]JavaInfo, error) {
	configs, err := LoadConfigs()
	if err != nil {
		return nil, err
	}
	return configs.Javas, nil
}

func SaveJavaLists(configs []JavaInfo) error {
	MCSCSConfigs, err := LoadConfigs()
	if err != nil {
		return err
	}
	MCSCSConfigs.Javas = configs
	err = SaveConfigs(MCSCSConfigs)
	if err != nil {
		return err
	}
	return nil
}

type DownloadInfo struct {
	Path string     `json:"path"`
	Info ServerInfo `json:"info"`
}

func LoadDownloadsLists() ([]DownloadInfo, error) {
	configs, err := LoadConfigs()
	if err != nil {
		return nil, err
	}
	return configs.Downloads, nil
}

func SaveDownloadsLists(configs []DownloadInfo) error {
	MCSCSConfigs, err := LoadConfigs()
	if err != nil {
		return err
	}
	MCSCSConfigs.Downloads = configs
	err = SaveConfigs(MCSCSConfigs)
	if err != nil {
		return err
	}
	return nil
}
