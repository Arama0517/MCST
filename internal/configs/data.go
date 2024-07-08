/*
 * Minecraft Server Tool(MCST) is a command-line utility making Minecraft server creation quick and easy for beginners.
 * Copyright (c) 2024-2024 Arama.
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

package configs

import (
	"os"
	"path/filepath"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	Configs         Config
	DefaultSettings = Settings{
		Aria2: Aria2{
			Enabled:                true,
			RetryWait:              2,
			Split:                  5,
			MaxConnectionPerServer: 5,
			MinSplitSize:           "5M",
			Options:                []string{},
		},
		AutoAcceptEULA: false,
		Language:       language.English.String(),
	}
)

var (
	rootDir      string
	ServersDir   string
	DownloadsDir string
	configsPath  string
)

type Core struct {
	ID         int    `yaml:"id"`          // 核心id
	URL        string `yaml:"url"`         // 下载地址(如果不是本地的话)
	FileName   string `yaml:"file_name"`   // 文件名
	FilePath   string `yaml:"file_path"`   // 文件路径
	ExtrasData any    `yaml:"extras_data"` // 其他数据
}

type Java struct {
	Path      string   `yaml:"path"`       // Java路径
	Args      []string `yaml:"args"`       // Java虚拟机参数
	MaxMemory uint64   `yaml:"max_memory"` // Java虚拟机最大堆内存
	MinMemory uint64   `yaml:"min_memory"` // Java虚拟机初始堆内存
	Encoding  string   `yaml:"encoding"`   // 编码
}
type Server struct {
	Name       string   `yaml:"name"`        // 服务器名称
	Java       Java     `yaml:"java"`        // Java
	ServerArgs []string `yaml:"server_args"` // Minecraft服务器参数
}

type Aria2 struct {
	Enabled                bool     `yaml:"enabled"`
	RetryWait              int      `yaml:"retry_wait"`
	Split                  int      `yaml:"split"`
	MaxConnectionPerServer int      `yaml:"max_connection_per_server"`
	MinSplitSize           string   `yaml:"min_split_size"`
	Options                []string `yaml:"options"`
}

type Settings struct {
	Aria2          Aria2  `yaml:"aria2"`
	AutoAcceptEULA bool   `yaml:"auto_accept_eula"`
	Language       string `yaml:"language"`
}

type Config struct {
	Cores    map[int]Core      `yaml:"cores"`   // 核心列表
	Servers  map[string]Server `yaml:"servers"` // 服务器列表, 如果服务器名称(key)为temp, CreatePage调用时会视为暂存配置而不是名为temp的服务器
	Settings Settings          `yaml:"settings"`
}

func InitData() error {
	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	rootDir = filepath.Join(UserHomeDir, ".config", "MCST")
	ServersDir = filepath.Join(rootDir, "servers")
	DownloadsDir = filepath.Join(rootDir, "downloads")
	configsPath = filepath.Join(rootDir, "configs.yaml")

	if err := os.MkdirAll(rootDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(ServersDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(DownloadsDir, 0o755); err != nil {
		return err
	}

	// 初始化

	if _, err := os.Stat(configsPath); os.IsNotExist(err) {
		data, err := yaml.Marshal(Config{
			Settings: DefaultSettings,
		})
		if err != nil {
			return err
		}
		if err := os.WriteFile(configsPath, data, 0o644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	file, err := os.ReadFile(configsPath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, &Configs); err != nil {
		return err
	}
	return nil
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(configsPath, data, 0o644)
}
