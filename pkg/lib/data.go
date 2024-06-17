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

package lib

import (
	"encoding/json"
	"os"
)

var Configs Config

var (
	DataDir      string
	ServersDir   string
	DownloadsDir string
	ConfigsPath  string
)

type Core struct {
	ID         int    // 核心id
	URL        string `json:"url"`         // 下载地址(如果不是本地的话)
	FileName   string `json:"file_name"`   // 文件名
	FilePath   string `json:"file_path"`   // 文件路径
	ExtrasData any    `json:"extras_data"` // 其他数据
}

type Java struct {
	Path     string   `json:"path"`     // Java路径
	Args     []string `json:"args"`     // Java虚拟机参数
	Xmx      uint64   `json:"xmx"`      // Java虚拟机最大堆内存
	Xms      uint64   `json:"xms"`      // Java虚拟机初始堆内存
	Encoding string   `json:"encoding"` // 编码
}
type Server struct {
	Name       string   `json:"name"`        // 服务器名称
	Java       Java     `json:"java"`        // Java
	ServerArgs []string `json:"server_args"` // Minecraft服务器参数
}

type Aria2c struct {
	Enabled                bool     `json:"enabled"`
	RetryWait              int      `json:"retry_wait"`
	Split                  int      `json:"split"`
	MaxConnectionPerServer int      `json:"max_connection_per_server"`
	MinSplitSize           string   `json:"min_split_size"`
	Option                 []string `json:"option"`
}

type Config struct {
	Cores          map[int]Core      `json:"cores"`            // 核心列表
	Servers        map[string]Server `json:"servers"`          // 服务器列表, 如果服务器名称(key)为temp, CreatePage调用时会视为暂存配置而不是名为temp的服务器
	Aria2c         Aria2c            `json:"aria2c"`           // aria2c配置
	AutoAcceptEULA bool              `json:"auto_accept_eula"` // 是否自动同意EULA
}

func (c *Config) Save() error {
	configs, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigsPath, configs, 0o644)
}
