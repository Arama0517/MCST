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

package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"

	api "github.com/Arama0517/MCST/internal/API"
	"github.com/Arama0517/MCST/pkg/lib"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

func newDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "download",
		Short:             "下载核心",
		Long:              "从远程或本地获取核心",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
	}
	cmd.AddCommand(newListCoreCmd(), newLocalCmd(), newRemoteCmd(), newFastMirrorCmd(), newPolarsCmd())
	return cmd
}

func newListCoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "list",
		Aliases:           []string{"fm"},
		Short:             "列出所有核心",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			var keys []int
			for _, v := range lib.Configs.Cores {
				keys = append(keys, v.ID)
			}
			sort.Ints(keys)
			for _, key := range keys {
				data := lib.Configs.Cores[key]
				log.WithFields(log.Fields{
					"URL":  data.URL,
					"文件名":  data.FileName,
					"文件路径": data.FilePath,
				}).Info(fmt.Sprint(key))
			}
			return nil
		},
	}
}

func newLocalCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:               "local",
		Aliases:           []string{"lo"},
		Short:             "本地",
		Long:              "从本地获取核心",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			if _, err := os.Stat(path); err != nil {
				return err
			}
			lib.Configs.Cores[len(lib.Configs.Cores)] = lib.Core{
				ID:       len(lib.Configs.Cores),
				URL:      "unknown",
				FileName: filepath.Base(path),
				FilePath: path,
			}
			return lib.Configs.Save()
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "核心的路径")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func newRemoteCmd() *cobra.Command {
	var URL string
	cmd := &cobra.Command{
		Use:               "remote",
		Aliases:           []string{"r"},
		Short:             "远程",
		Long:              "从指定的URL下载核心",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			URL, err := url.Parse(URL)
			if err != nil {
				return err
			}
			path, err := lib.NewDownloader(*URL).Download()
			if err != nil {
				return err
			}
			lib.Configs.Cores[len(lib.Configs.Cores)] = lib.Core{
				ID:       len(lib.Configs.Cores),
				URL:      URL.String(),
				FileName: filepath.Base(path),
				FilePath: path,
			}
			return lib.Configs.Save()
		},
	}
	cmd.Flags().StringVarP(&URL, "url", "u", "", "核心的URL")
	_ = cmd.MarkFlagRequired("url")
	return cmd
}

type fastMirrorCmdFlags struct {
	core             string
	minecraftVersion string
	buildVersion     string
}

func newFastMirrorCmd() *cobra.Command {
	flags := fastMirrorCmdFlags{}
	cmd := &cobra.Command{
		Use:               "fastmirror",
		Aliases:           []string{"fm"},
		Short:             "从无极镜像下载核心",
		Long:              "无极镜像: <https://www.fastmirror.net/>\n请使用 'MCST download fastmirror list' 获取所需的参数",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			downloader := api.GetFastMirrorDownloader(flags.core, flags.minecraftVersion, flags.buildVersion)
			path, err := downloader.Download()
			if err != nil {
				return err
			}
			lib.Configs.Cores[len(lib.Configs.Cores)] = lib.Core{
				ID:       len(lib.Configs.Cores),
				URL:      downloader.URL.String(),
				FileName: downloader.FileName,
				FilePath: path,
				ExtrasData: map[string]any{
					"core":          flags.core,
					"mc_version":    flags.minecraftVersion,
					"build_version": flags.buildVersion,
				},
			}
			return lib.Configs.Save()
		},
	}
	cmd.AddCommand(newListFastMirrorCmd())
	cmd.Flags().StringVarP(&flags.core, "core", "c", "", "核心的名称, 例如: \"Mohist\"")
	cmd.Flags().StringVarP(&flags.minecraftVersion, "mc_version", "m", "", "核心的Minecraft版本, 例如: \"1.20.1\"")
	cmd.Flags().StringVarP(&flags.buildVersion, "build_version", "b", "", "核心的构建版本, 例如: \"build738\"")
	_ = cmd.MarkFlagRequired("core")
	_ = cmd.MarkFlagRequired("mc_version")
	_ = cmd.MarkFlagRequired("build_version")
	return cmd
}

func newListFastMirrorCmd() *cobra.Command {
	flags := fastMirrorCmdFlags{}
	cmd := &cobra.Command{
		Use:   "list",
		Short: "获取无极镜像的核心信息",
		Long: `使用此命令时拥有三种情况:
1. 不使用任何参数: 输出所有核心的名称
2. 使用 '--core': 输出此核心的所有Minecraft版本
3. 使用 '--core' 和 '--mc_version': 返回此核心的所有构建版本`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			fastMirror, err := api.GetFastMirrorData()
			if err != nil {
				return err
			}
			switch cmdFlags := cmd.Flags(); {
			case !cmdFlags.Changed("core") && !cmdFlags.Changed("mc_version"):
				keys := make([]string, 0, len(fastMirror))
				for k := range fastMirror {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, key := range keys {
					data := fastMirror[key]
					log.WithFields(log.Fields{
						"推荐": data.Recommend,
						"标签": data.Tag,
						"主页": data.Homepage,
					}).Info(data.Name)
				}
			case cmdFlags.Changed("core") && !cmdFlags.Changed("mc_version"):
				data := fastMirror[flags.core].MinecraftVersions
				sort.Strings(data)
				for _, version := range data {
					log.Info(version)
				}
			case cmdFlags.Changed("core") && cmdFlags.Changed("mc_version"):
				fastMirror, err := api.GetFastMirrorBuildsData(flags.core, flags.minecraftVersion)
				if err != nil {
					return err
				}
				keys := make([]string, 0, len(fastMirror))
				for k := range fastMirror {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, key := range keys {
					data := fastMirror[key]
					updateTime, err := time.Parse("2006-01-02T15:04:05", data.UpdateTime)
					if err != nil {
						return err
					}
					log.WithFields(log.Fields{
						"SHA1": data.Sha1,
						"更新时间": updateTime.Format("2006-01-02 15:04:05"),
					}).Info(data.CoreVersion)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.core, "core", "c", "", "核心的名称, 例如: \"Mohist\"")
	cmd.Flags().StringVarP(&flags.minecraftVersion, "mc_version", "m", "", "核心的Minecraft版本, 例如: \"1.20.1\"")
	return cmd
}

func newPolarsCmd() *cobra.Command {
	var typeID, coreID int
	cmd := &cobra.Command{
		Use:               "polars",
		Short:             "从极星云镜像下载核心",
		Long:              "极星云镜像: <https://mirror.polars.cc/>",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			polars, err := api.GetPolarsCoresData(typeID)
			if err != nil {
				return err
			}
			URL, err := url.Parse(polars[coreID].DownloadURL)
			if err != nil {
				return err
			}
			downloader := lib.NewDownloader(*URL)
			path, err := downloader.Download()
			if err != nil {
				return err
			}
			lib.Configs.Cores[len(lib.Configs.Cores)] = lib.Core{
				ID:       len(lib.Configs.Cores),
				URL:      URL.String(),
				FileName: downloader.FileName,
				FilePath: path,
				ExtrasData: map[string]int{
					"type_id": typeID,
					"core_id": coreID,
				},
			}
			return lib.Configs.Save()
		},
	}
	cmd.AddCommand(newListPolarsCmd())
	cmd.Flags().IntVar(&typeID, "type_id", 0, "核心类型的ID")
	cmd.Flags().IntVar(&coreID, "core_id", 0, "核心ID")
	return cmd
}

func newListPolarsCmd() *cobra.Command {
	var id int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "获取极星云镜像的核心信息",
		Long: `使用此命令时拥有两种情况:
1. 不使用任何参数: 输出所有核心类型的ID
2. 使用 '--id': 输出此核心的所有的核心ID(用于下载)`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cmd.Flags().Changed("id") {
				polars, err := api.GetPolarsCoresData(id)
				if err != nil {
					return err
				}
				keys := make([]int, 0, len(polars))
				for k := range polars {
					keys = append(keys, k)
				}
				sort.Ints(keys)
				for _, k := range keys {
					data := polars[k]
					log.WithFields(log.Fields{
						"名称":   data.Name,
						"下载链接": data.DownloadURL,
					}).Info(fmt.Sprint(data.ID))
				}
			} else {
				polars, err := api.GetPolarsData()
				if err != nil {
					return err
				}
				keys := make([]string, 0, len(polars))
				for k := range polars {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					data := polars[k]
					log.WithFields(log.Fields{
						"名称": data.Name,
						"描述": data.Description,
					}).Info(fmt.Sprint(data.ID))
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, "类型ID")
	return cmd
}
