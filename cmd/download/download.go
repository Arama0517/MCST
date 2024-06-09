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
	"os"
	"path/filepath"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

func NewDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "下载核心",
		Long:  "从远程或本地获取核心",
	}
	cmd.AddCommand(newListCmd(), newLocalCmd(), newRemoteCmd(), newFastMirrorCmd(), newPolarsCmd())
	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"fm"},
		Short:   "列出所有核心",
		RunE: func(*cobra.Command, []string) error {
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			for _, core := range configs.Cores {
				log.WithFields(log.Fields{
					"URL":  core.URL,
					"文件名":  core.FileName,
					"文件路径": core.FilePath,
					"其他数据": core.ExtrasData.(log.Fields),
				}).Info(fmt.Sprint(core.ID))
			}
			return nil
		},
	}
}

func newLocalCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:     "local",
		Aliases: []string{"lo"},
		Short:   "本地",
		Long:    "从本地获取核心",
		RunE: func(*cobra.Command, []string) error {
			if _, err := os.Stat(path); err != nil {
				return err
			}
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			configs.Cores = append(configs.Cores, lib.Core{
				URL:      "unknown",
				FileName: filepath.Base(path),
				FilePath: path,
				ID:       len(configs.Cores) + 1,
			})
			return configs.Save()
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "核心的路径")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func newRemoteCmd() *cobra.Command {
	var URL string
	cmd := &cobra.Command{
		Use:     "remote",
		Aliases: []string{"r"},
		Short:   "远程",
		Long:    "从指定的URL下载核心",
		RunE: func(*cobra.Command, []string) error {
			URL, err := url.Parse(URL)
			if err != nil {
				return err
			}
			path, err := lib.NewDownloader(*URL).Download()
			if err != nil {
				return err
			}
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			configs.Cores = append(configs.Cores, lib.Core{
				URL:      URL.String(),
				FileName: filepath.Base(path),
				FilePath: path,
				ID:       len(configs.Cores) + 1,
			})
			return configs.Save()
		},
	}
	cmd.Flags().StringVarP(&URL, "url", "u", "", "核心的URL")
	_ = cmd.MarkFlagRequired("url")
	return cmd
}
