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
	"net/url"
	"sort"

	api "github.com/Arama-Vanarana/MCServerTool/internal/API"
	lib2 "github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

func newPolarsCmd() *cobra.Command {
	var typeID, coreID int
	cmd := &cobra.Command{
		Use:   "polars",
		Short: "从极星云镜像下载核心",
		Long:  "极星云镜像: <https://mirror.polars.cc/>",
		RunE: func(*cobra.Command, []string) error {
			polars, err := api.GetPolarsCoresDatas(typeID)
			if err != nil {
				return err
			}
			URL, err := url.Parse(polars[coreID].DownloadURL)
			if err != nil {
				return err
			}
			downloader := lib2.NewDownloader(*URL)
			path, err := downloader.Download()
			if err != nil {
				return err
			}
			configs, err := lib2.LoadConfigs()
			if err != nil {
				return err
			}
			configs.Cores = append(configs.Cores, lib2.Core{
				URL:      URL.String(),
				FileName: downloader.FileName,
				FilePath: path,
				ExtrasData: map[string]int{
					"type_id": typeID,
					"core_id": coreID,
				},
			})
			return configs.Save()
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cmd.Flags().Changed("id") {
				polars, err := api.GetPolarsCoresDatas(id)
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
