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
	"net/url"
	"os"
	"path/filepath"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

func newDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "下载核心",
		Long:  "从远程或本地获取核心",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "列出所有核心",
		RunE: func(cmd *cobra.Command, _ []string) error {
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			for index, core := range configs.Cores {
				log.WithField("info", core).WithField("index", index).Info(core.FileName)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd, newLocalCmd(), newRemoteCmd())
	return cmd
}

func newLocalCmd() *cobra.Command {
	var path string
	cmd := &cobra.Command{
		Use:   "local",
		Short: "本地",
		Long:  "从本地获取核心",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if _, err := os.Stat(path); err != nil {
				return err
			}
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			configs.Cores = append(configs.Cores, lib.Core{
				FileName: filepath.Base(path),
				FilePath: path,
			})
			if err := configs.Save(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", "核心的路径")
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func newRemoteCmd() *cobra.Command {
	var URL string
	cmd := &cobra.Command{
		Use:   "remote",
		Short: "远程",
		Long:  "从指定的URL下载核心",
		RunE: func(cmd *cobra.Command, _ []string) error {
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
				URL:      *URL,
				FileName: filepath.Base(path),
				FilePath: path,
			})
			if err := configs.Save(); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&URL, "url", "u", "", "核心的URL")
	_ = cmd.MarkFlagRequired("url")
	return cmd
}
