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

package download

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Arama0517/MCST/internal/build"
	"github.com/Arama0517/MCST/internal/configs"
	"github.com/go-cmd/cmd"
	"github.com/siku2/arigo"
)

const StatusComplete arigo.DownloadStatus = "complete"

func (d *Downloader) aria2Download() (string, error) {
	// 启动Aria2
	aria2Cmd := cmd.NewCmd("aria2c")
	aria2Cmd.Args = append(aria2Cmd.Args, configs.Configs.Settings.Aria2.Options...)
	aria2Cmd.Args = append(aria2Cmd.Args,
		fmt.Sprintf("--dir=%s", configs.DownloadsDir),
		fmt.Sprintf("--user-agent=MCST/%s", build.Version.GitVersion),
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		fmt.Sprintf("--retry-wait=%d", configs.Configs.Settings.Aria2.RetryWait),
		fmt.Sprintf("--split=%d", configs.Configs.Settings.Aria2.Split),
		fmt.Sprintf("--max-connection-per-server=%d", configs.Configs.Settings.Aria2.MaxConnectionPerServer),
		fmt.Sprintf("--min-split-size=%s", configs.Configs.Settings.Aria2.MinSplitSize),
		"--enable-rpc",
		"--rpc-listen-all",
		"--rpc-listen-port=6800",
		"--rpc-secret=MCST",
		"--console-log-level=notice",
		"--no-conf=true",
		"--follow-metalink=true",
		"--metalink-preferred-protocol=https",
		"--min-tls-version=TLSv1.2",
		fmt.Sprintf("--stop-with-process=%d", os.Getpid()),
		"--continue",
		"--summary-interval=0",
		"--auto-save-interval=1",
	)
	aria2Cmd.Start()

	// 检测是否启动JSONRPC
	for {
		status := aria2Cmd.Status()
		n := len(status.Stdout)
		if strings.Contains(status.Stdout[n-1], "IPV4 RPC") {
			break
		}
	}

	// 连接Aria2
	client, err := arigo.Dial("ws://localhost:6800/jsonrpc", "MCST")
	if err != nil {
		return "", err
	}

	// 获取GID
	gid, err := client.AddURI(arigo.URIs(d.URL), nil)
	if err != nil {
		return "", err
	}

	// 更新进度条
	for {
		status, err := gid.TellStatus("status", "completedLength", "connections")
		if err != nil {
			return "", err
		}
		if status.Status == StatusComplete {
			files, err := gid.GetFiles()
			if err != nil {
				return "", err
			}
			if err = d.bar.Finish(); err != nil {
				return "", err
			}
			return files[0].Path, nil
		}
		d.bar.Describe(fmt.Sprintf("[green]已连接至%d个服务器[reset] [cyan]下载中...[reset]", status.Connections))
		if err = d.bar.Set64(int64(status.CompletedLength)); err != nil {
			return "", err
		}
		time.Sleep(50 * time.Millisecond)
	}
}
