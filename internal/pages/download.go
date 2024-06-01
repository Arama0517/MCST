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
	"fmt"
	"net/url"
	"path/filepath"

	api "github.com/Arama-Vanarana/MCServerTool/pkg/API"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

var Download = cli.Command{
	Name:  "download",
	Usage: "下载核心",
	Subcommands: []*cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "查看已下载的核心",
			Action: func(_ *cli.Context) error {
				configs, err := lib.LoadConfigs()
				if err != nil {
					return err
				}
				for i, Core := range configs.Cores {
					fmt.Printf("%s(%d): %s\n", Core.FileName, i, Core.FilePath)
				}
				return nil
			},
		},
		{
			Name:  "local",
			Usage: "使用本地核心",
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name:     "path",
					Aliases:  []string{"p"},
					Usage:    "本地核心路径",
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) error {
				configs, err := lib.LoadConfigs()
				if err != nil {
					return err
				}
				configs.Cores = append(configs.Cores, lib.Core{
					FileName: ctx.Path("path"),
				})
				if err := configs.Save(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:    "remote",
			Aliases: []string{"r"},
			Usage:   "从指定的URL下载核心",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "url",
					Aliases:  []string{"u"},
					Usage:    "下载核心的URL",
					Required: true,
				},
			},
			Action: func(ctx *cli.Context) error {
				configs, err := lib.LoadConfigs()
				if err != nil {
					return err
				}
				url, err := url.Parse(ctx.String("url"))
				if err != nil {
					return err
				}
				path, err := (&lib.Downloader{
					URL: *url,
				}).Download()
				if err != nil {
					return err
				}
				configs.Cores = append(configs.Cores, lib.Core{
					URL:      *url,
					FileName: filepath.Base(path),
					FilePath: path,
				})
				if err := configs.Save(); err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:    "FastMirror",
			Aliases: []string{"fm"},
			Usage:   "从无极镜像(https://www.fastmirror.net)下载核心, 如果不使用 '-l' 或 '--list' 参数就会下载指定的版本(必须含有 '-c' , '-m' 和 '-b' 参数)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "core",
					Aliases: []string{"c"},
					Usage:   "服务器核心, 例如: Mohist",
				},
				&cli.StringFlag{
					Name:    "mc_version",
					Aliases: []string{"m"},
					Usage:   "Minecraft版本, 例如: 1.20.1",
				},
				&cli.StringFlag{
					Name:    "build_version",
					Aliases: []string{"b"},
					Usage:   "构建版本, 例如: build524",
				},
				&cli.BoolFlag{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "列出无极镜像可用版本, 例如 '-c Mohist -l' 参数就会输出Mohist的可用版本, 使用 '-c Mohist -m 1.20.1 -l' 参数就会返回Mohist 1.20.1的可用构建版本",
				},
			},
			Action: fastMirror,
		},
		{
			Name:    "Polars",
			Aliases: []string{"pl"},
			Usage:   "从极星云镜像(https://mirror.polars.cc)下载核心, 如果不使用 '-l' 或 '--list' 参数就会下载指定的版本(必须含有 '--id' 参数) 不推荐使用(因为核心更新时间较落后)",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "type_id",
					Aliases: []string{"ti"},
					Usage:   "服务器类型ID, 例如: 1",
				},
				&cli.IntFlag{
					Name:    "core_id",
					Aliases: []string{"ci"},
					Usage:   "服务器核心ID, 例如: 1",
				},
				&cli.BoolFlag{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "列出极星云镜像可用版本, 例如不带任何参数就会输出所有可用的核心和他的ID, 使用 '--id' 参数就会输出指定核心的可用版本",
				},
			},
			Action: polars,
		},
	},
}

func fastMirror(ctx *cli.Context) error {
	core := ctx.String("core")
	minecraftVersion := ctx.String("mc_version")
	buildVersion := ctx.String("build_version")
	list := ctx.Bool("list")
	fastMirror, err := api.GetFastMirrorDatas()
	if err != nil {
		return err
	}
	if list {
		switch {
		case core != "" && minecraftVersion != "":
			fastMirrorBuilds, err := api.GetFastMirrorBuildsDatas(core, minecraftVersion)
			if err != nil {
				return err
			}
			for _, data := range fastMirrorBuilds {
				fmt.Printf("%s: 更新时间: %s, SHA1: %s\n", data.CoreVersion, data.UpdateTime, data.Sha1)
			}
		case core != "":
			for _, data := range fastMirror[core].MinecraftVersions {
				fmt.Println(data)
			}
		default:
			return errors.New("没有这个用法")
		}
	} else {
		if core == "" || minecraftVersion == "" || buildVersion == "" {
			return errors.New("缺少必要参数")
		}
		path, err := api.DownloadFastMirrorServer(core, minecraftVersion, buildVersion)
		if err != nil {
			return err
		}
		configs, err := lib.LoadConfigs()
		if err != nil {
			return err
		}
		configs.Cores = append(configs.Cores, lib.Core{
			FileName: filepath.Base(path),
			FilePath: path,
			ExtrasData: map[string]string{
				"core":          core,
				"mc_version":    minecraftVersion,
				"build_version": buildVersion,
			},
		})
		if err := configs.Save(); err != nil {
			return err
		}
	}
	return nil
}

func polars(ctx *cli.Context) error {
	typeID := ctx.Int("type_id")
	coreID := ctx.Int("core_id")
	list := ctx.Bool("list")
	polars, err := api.GetPolarsData()
	if err != nil {
		return err
	}
	if list {
		switch {
		case typeID == 0 && coreID == 0:
			for _, data := range polars {
				fmt.Printf("%s(%d): %s\n", data.Name, data.ID, data.Description)
			}
		case typeID != 0 && coreID == 0:
			data, err := api.GetPolarsCoresDatas(typeID)
			if err != nil {
				return err
			}
			for _, core := range data {
				fmt.Printf("%s(%d): %s\n", core.Name, core.ID, core.DownloadURL)
			}
		default:
			return errors.New("没有这个用法")
		}
	} else {
		if typeID == 0 || coreID == 0 {
			return errors.New("缺少必要参数")
		}
		data, err := api.GetPolarsCoresDatas(typeID)
		if err != nil {
			return err
		}
		path, err := api.DownloadPolarsServer(data[coreID].DownloadURL)
		if err != nil {
			return err
		}
		configs, err := lib.LoadConfigs()
		if err != nil {
			return err
		}
		URL, err := url.Parse(data[coreID].DownloadURL)
		if err != nil {
			return err
		}
		configs.Cores = append(configs.Cores, lib.Core{
			URL:      *URL,
			FileName: filepath.Base(path),
			FilePath: path,
			ExtrasData: map[string]int{
				"type_id": typeID,
				"core_id": coreID,
			},
		})
		if err := configs.Save(); err != nil {
			return err
		}
	}
	return nil
}
