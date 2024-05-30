/*
 * Minecraft Server Tool(MST) is a command-line utility making Minecraft server creation quick and easy for beginners.
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

package main

import (
	"fmt"
	"os"

	"github.com/Arama-Vanarana/MCServerTool/internal/pages"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

func init() {
	lib.Init()
}

func main() {
	app := cli.App{
		Name:    "MST",
		Usage:   "Minecraft Server Tool",
		Version: lib.Version,
		Authors: []*cli.Author{
			{
				Name:  "Arama",
				Email: "3584075812@qq.com",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "创建服务器",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "name",
						Aliases:  []string{"n"},
						Usage:    "服务器名称",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "Xms",
						Aliases: []string{"m"},
						Usage:   "Xms, Java虚拟机初始堆内存(可用单位: B, KB, KiB, MB, MiB, GB, GiB, TB, TiB), 默认为1G",
						Value:   "1G",
					},
					&cli.StringFlag{
						Name:    "Xmx",
						Aliases: []string{"x"},
						Usage:   "Xmx, Java虚拟机最大堆内存(可用单位: B, KB, KiB, MB, MiB, GB, GiB, TB, TiB); 默认为1G",
						Value:   "1G",
					},
					&cli.BoolFlag{
						Name:  "gbk",
						Usage: "是否使用GBK编码, 默认为false",
						Value: false,
					},
					&cli.PathFlag{
						Name:     "java",
						Aliases:  []string{"j"},
						Usage:    "Java路径",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "jvm_args",
						Aliases: []string{"a"},
						Usage:   "Java虚拟机的参数, 例如: -a '-arg1 -arg2', 用空格分隔参数",
						Value:   "-Dlog4j2.formatMsgNoLookups=true",
					},
					&cli.IntFlag{
						Name:     "core",
						Aliases:  []string{"c"},
						Usage:    "服务器核心, 从download命令下载, 使用downloads --list查看已下载的核心",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "server_args",
						Aliases: []string{"s"},
						Usage:   "Minecraft服务器特有参数, 例如: -s '--nogui --arg2', 用空格分隔参数",
						Value:   "--nogui",
					},
				},
				Action: pages.Create,
			},
			{
				Name:  "download",
				Usage: "下载核心",
				Subcommands: []*cli.Command{
					{
						Name:   "list",
						Usage:  "查看已下载的核心",
						Action: pages.ListCores,
					},
					{
						Name:    "local",
						Aliases: []string{"l"},
						Usage:   "使用本地核心",
						Flags: []cli.Flag{
							&cli.PathFlag{
								Name:     "path",
								Aliases:  []string{"p"},
								Usage:    "本地核心路径",
								Required: true,
							},
						},
						Action: pages.Local,
					},
					{
						Name:    "remote",
						Aliases: []string{"r"},
						Usage:   "从指定的URL下载核心",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "URL",
								Aliases:  []string{"u"},
								Usage:    "下载核心的URL",
								Required: true,
							},
						},
						Action: pages.Remote,
					},
					{
						Name:    "FastMirror",
						Aliases: []string{"fm"},
						Usage:   "从无极镜像(https://www.fastmirror.net)下载核心, 如果不使用 '-l' 或 '--list' 参数就会下载指定的版本(必须含有 '-c' , '-m' 和 '-b' 参数)",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "Core",
								Aliases: []string{"c"},
								Usage:   "服务器核心, 例如: Mohist",
							},
							&cli.StringFlag{
								Name:    "MinecraftVersion",
								Aliases: []string{"m"},
								Usage:   "Minecraft版本, 例如: 1.20.1",
							},
							&cli.StringFlag{
								Name:    "BuildVersion",
								Aliases: []string{"b"},
								Usage:   "构建版本, 例如: build524",
							},
							&cli.BoolFlag{
								Name:    "list",
								Aliases: []string{"l"},
								Usage:   "列出无极镜像可用版本, 例如 '-c Mohist -l' 参数就会输出Mohist的可用版本, 使用 '-c Mohist -m 1.20.1 -l' 参数就会返回Mohist 1.20.1的可用构建版本",
							},
						},
						Action: pages.FastMirror,
					},
					{
						Name:    "Polars",
						Aliases: []string{"pl"},
						Usage:   "从极星云镜像(https://mirror.polars.cc)下载核心, 如果不使用 '-l' 或 '--list' 参数就会下载指定的版本(必须含有 '--id' 参数) 不推荐使用(因为核心更新时间较落后)",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "TypeID",
								Aliases: []string{"ti"},
								Usage:   "服务器类型ID, 例如: 1",
							},
							&cli.IntFlag{
								Name:    "CoreID",
								Aliases: []string{"ci"},
								Usage:   "服务器核心ID, 例如: 1",
							},
							&cli.BoolFlag{
								Name:    "list",
								Aliases: []string{"l"},
								Usage:   "列出极星云镜像可用版本, 例如不带任何参数就会输出所有可用的核心和他的ID, 使用 '--id' 参数就会输出指定核心的可用版本",
							},
						},
						Action: pages.Polars,
					},
				},
			},
			{
				Name:   "list",
				Usage:  "列出服务器",
				Action: pages.ListServers,
			},
			{
				Name:  "start",
				Usage: "启动服务器",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "server",
						Aliases:  []string{"s"},
						Usage:    "要启动的服务器名称",
						Required: true,
					},
				},
				Action: pages.Start,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
