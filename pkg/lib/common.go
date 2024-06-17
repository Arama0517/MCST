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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	goversion "github.com/caarlos0/go-version"
)

// Init 一键全部初始化(按顺序)
func Init(v goversion.Info) error {
	// 设置版本
	version = v

	UserHomeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	DataDir = filepath.Join(UserHomeDir, ".config", "MCST")
	ServersDir = filepath.Join(DataDir, "servers")
	DownloadsDir = filepath.Join(DataDir, "downloads")
	ConfigsPath = filepath.Join(DataDir, "configs.json")

	// 初始化
	if _, err := os.Stat(DataDir); os.IsNotExist(err) {
		if err := os.MkdirAll(DataDir, 0o755); err != nil {
			return err
		}
		if err := os.MkdirAll(ServersDir, 0o755); err != nil {
			return err
		}
		if err := os.MkdirAll(DownloadsDir, 0o755); err != nil {
			return err
		}
		jsonData, err := json.MarshalIndent(Config{
			Cores:   map[int]Core{},
			Servers: map[string]Server{},
			Aria2c: Aria2c{
				Enabled:                true,
				RetryWait:              2,
				Split:                  5,
				MaxConnectionPerServer: 5,
				MinSplitSize:           "5M",
			},
			AutoAcceptEULA: false,
		}, "", "    ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(ConfigsPath, jsonData, 0o644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	file, err := os.ReadFile(ConfigsPath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(file, &Configs); err != nil {
		return err
	}

	return nil
}

var version goversion.Info

// Request 请求URL, 返回响应; 运行成功后请添加`defer resp.Body.Close()`到你的代码内
func Request(url url.URL, method string, header map[string]string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("User-Agent", fmt.Sprintf("MCST/%s ", version.GitVersion))
	for key, value := range header {
		request.Header.Set(key, value)
	}
	return http.DefaultClient.Do(request)
}
