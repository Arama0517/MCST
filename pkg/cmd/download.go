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
	"os"
	"path/filepath"
	"sort"
	"time"

	api "github.com/Arama0517/MCST/internal/API"
	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/download"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

func newDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "download",
		Short:             locale.GetLocaleMessage("download.short"),
		Long:              locale.GetLocaleMessage("download.long"),
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
		Short:             locale.GetLocaleMessage("download.list"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			var keys []int
			for _, v := range configs.Configs.Cores {
				keys = append(keys, v.ID)
			}
			sort.Ints(keys)
			for _, key := range keys {
				data := configs.Configs.Cores[key]
				log.WithFields(log.Fields{
					locale.GetLocaleMessage("download.list.output.url"):      data.URL,
					locale.GetLocaleMessage("download.list.output.filename"): data.FileName,
					locale.GetLocaleMessage("download.list.output.filepath"): data.FilePath,
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
		Short:             locale.GetLocaleMessage("download.local.short"),
		Long:              locale.GetLocaleMessage("download.local.long"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			if _, err := os.Stat(path); err != nil {
				return err
			}
			configs.Configs.Cores[len(configs.Configs.Cores)] = configs.Core{
				ID:       len(configs.Configs.Cores),
				URL:      "unknown",
				FileName: filepath.Base(path),
				FilePath: path,
			}
			return configs.Configs.Save()
		},
	}
	cmd.Flags().StringVarP(&path, "path", "p", "", locale.GetLocaleMessage("download.local.flags.path"))
	_ = cmd.MarkFlagRequired("path")
	return cmd
}

func newRemoteCmd() *cobra.Command {
	var URL string
	cmd := &cobra.Command{
		Use:               "remote",
		Short:             locale.GetLocaleMessage("download.remote.short"),
		Long:              locale.GetLocaleMessage("download.remote.long"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			downloader := download.NewDownloader(URL)
			path, err := downloader.Download()
			if err != nil {
				return err
			}
			configs.Configs.Cores[len(configs.Configs.Cores)] = configs.Core{
				ID:       len(configs.Configs.Cores),
				URL:      URL,
				FileName: downloader.FileName,
				FilePath: path,
			}
			return configs.Configs.Save()
		},
	}
	cmd.Flags().StringVarP(&URL, "url", "u", "", locale.GetLocaleMessage("download.remote.flags.url"))
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
		Short:             locale.GetLocaleMessage("download.fastmirror.short"),
		Long:              locale.GetLocaleMessage("download.fastmirror.long"),
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
			configs.Configs.Cores[len(configs.Configs.Cores)] = configs.Core{
				ID:       len(configs.Configs.Cores),
				URL:      downloader.URL,
				FileName: downloader.FileName,
				FilePath: path,
				ExtrasData: map[string]any{
					"core":          flags.core,
					"mc_version":    flags.minecraftVersion,
					"build_version": flags.buildVersion,
				},
			}
			return configs.Configs.Save()
		},
	}
	cmd.AddCommand(newListFastMirrorCmd())
	cmd.Flags().StringVarP(&flags.core, "core", "c", "", locale.GetLocaleMessage("download.fastmirror.flags.core"))
	cmd.Flags().StringVarP(&flags.minecraftVersion, "mc_version", "m", "", locale.GetLocaleMessage("download.fastmirror.flags.mc_version"))
	cmd.Flags().StringVarP(&flags.buildVersion, "build_version", "b", "", locale.GetLocaleMessage("download.fastmirror.flags.build_version"))
	_ = cmd.MarkFlagRequired("core")
	_ = cmd.MarkFlagRequired("mc_version")
	_ = cmd.MarkFlagRequired("build_version")
	return cmd
}

func newListFastMirrorCmd() *cobra.Command {
	flags := fastMirrorCmdFlags{}
	cmd := &cobra.Command{
		Use:               "list",
		Short:             locale.GetLocaleMessage("download.fastmirror.list.short"),
		Long:              locale.GetLocaleMessage("download.fastmirror.list.long"),
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
						locale.GetLocaleMessage("download.fastmirror.list.output.recommend"): data.Recommend,
						locale.GetLocaleMessage("download.fastmirror.list.output.tag"):       data.Tag,
						locale.GetLocaleMessage("download.fastmirror.list.output.homepage"):  data.Homepage,
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
						locale.GetLocaleMessage("download.fastmirror.list.output.sha1"):        data.Sha1,
						locale.GetLocaleMessage("download.fastmirror.list.output.update_time"): updateTime.Format("2006-01-02 15:04:05"),
					}).Info(data.CoreVersion)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.core, "core", "c", "", locale.GetLocaleMessage("download.fastmirror.flags.core"))
	cmd.Flags().StringVarP(&flags.minecraftVersion, "mc_version", "m", "", locale.GetLocaleMessage("download.fastmirror.flags.mc_version"))
	return cmd
}

func newPolarsCmd() *cobra.Command {
	var typeID, coreID int
	cmd := &cobra.Command{
		Use:               "polars",
		Short:             locale.GetLocaleMessage("download.polars.short"),
		Long:              locale.GetLocaleMessage("download.polars.long"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			polars, err := api.GetPolarsCoresData(typeID)
			if err != nil {
				return err
			}
			downloader := download.NewDownloader(polars[coreID].DownloadURL)
			path, err := downloader.Download()
			if err != nil {
				return err
			}
			configs.Configs.Cores[len(configs.Configs.Cores)] = configs.Core{
				ID:       len(configs.Configs.Cores),
				URL:      polars[coreID].DownloadURL,
				FileName: downloader.FileName,
				FilePath: path,
				ExtrasData: map[string]int{
					"type_id": typeID,
					"core_id": coreID,
				},
			}
			return configs.Configs.Save()
		},
	}
	cmd.AddCommand(newListPolarsCmd())
	cmd.Flags().IntVar(&typeID, "type_id", 0, locale.GetLocaleMessage("download.polars.flags.type_id"))
	cmd.Flags().IntVar(&coreID, "core_id", 0, locale.GetLocaleMessage("download.polars.flags.core_id"))
	return cmd
}

func newListPolarsCmd() *cobra.Command {
	var id int
	cmd := &cobra.Command{
		Use:               "list",
		Short:             locale.GetLocaleMessage("download.polars.list.short"),
		Long:              locale.GetLocaleMessage("download.polars.list.long"),
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
						locale.GetLocaleMessage("download.polars.output.file_name"):    data.Name,
						locale.GetLocaleMessage("download.polars.output.download_url"): data.DownloadURL,
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
						locale.GetLocaleMessage("download.polars.output.name"):        data.Name,
						locale.GetLocaleMessage("download.polars.output.description"): data.Description,
					}).Info(fmt.Sprint(data.ID))
				}
			}
			return nil
		},
	}
	cmd.Flags().IntVar(&id, "id", 0, locale.GetLocaleMessage("download.polars.flags.type_id"))
	return cmd
}
