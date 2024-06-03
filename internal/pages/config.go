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

package pages

import (
	"errors"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/urfave/cli/v2"
)

var Config = cli.Command{
	Name:  "config",
	Usage: "配置服务器",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "server",
			Aliases:  []string{"s"},
			Usage:    "要配置的服务器名称",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "服务器名称",
		},
		&cli.StringFlag{
			Name:    "Xms",
			Aliases: []string{"min"},
			Usage:   "Xms, Java虚拟机初始堆内存(可用单位: B, KB, KiB, MB, MiB, GB, GiB, TB, TiB)",
		},
		&cli.StringFlag{
			Name:    "Xmx",
			Aliases: []string{"max"},
			Usage:   "Xmx, Java虚拟机最大堆内存(可用单位: B, KB, KiB, MB, MiB, GB, GiB, TB, TiB)",
		},
		&cli.StringFlag{
			Name:    "encoding",
			Aliases: []string{"e"},
			Usage:   "使用的编码格式, 常见的编码格式: UTF-8, GBK, ASCII",
		},
		&cli.PathFlag{
			Name:    "java",
			Aliases: []string{"j"},
			Usage:   "Java路径",
		},
		&cli.StringSliceFlag{
			Name:    "jvm_args",
			Aliases: []string{"ja"},
			Usage:   "Java虚拟机的参数",
		},
		&cli.StringSliceFlag{
			Name:    "server_args",
			Aliases: []string{"sa"},
			Usage:   "Minecraft服务器特有参数",
		},
	},
	Action: func(context *cli.Context) error {
		configs, err := lib.LoadConfigs()
		if err != nil {
			return err
		}
		server, exists := configs.Servers[context.String("server")]
		if !exists {
			return errors.New("服务器不存在")
		}
		name := context.String("name")
		if name != "" {
			server.Name = name
		}
		Xms, err := toBytes(context.String("Xms"))
		if err != nil {
			return err
		}
		Xmx, err := toBytes(context.String("Xmx"))
		if err != nil {
			return err
		}
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			return err
		}
		switch {
		case Xms > server.Java.Xmx, Xms > Xmx:
			return errors.New("Xms不能大于Xmx")
		case Xms < 1024*1024, Xmx < 1024*1024:
			return errors.New("Xms和Xmx必须大于1MiB")
		case Xms > memInfo.Total, Xmx > memInfo.Total:
			return errors.New("Xms和Xmx不能大于系统内存")
		}
		server.Java.Xms = Xms
		server.Java.Xmx = Xmx
		encoding := context.String("encoding")
		if encoding != "" {
			server.Java.Encoding = encoding
		}
		java := context.Path("java")
		if java != "" {
			server.Java.Path = java
		}
		jvmArgs := context.StringSlice("jvm_args")
		if len(jvmArgs) > 0 {
			server.Java.Args = jvmArgs
		}
		serverArgs := context.StringSlice("server_args")
		if len(serverArgs) > 0 {
			server.ServerArgs = serverArgs
		}
		return nil
	},
}
