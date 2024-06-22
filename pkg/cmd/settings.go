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

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

func newSettingsCmd() *cobra.Command {
	var flags configs.Config
	var Language string
	cmd := &cobra.Command{
		Use:               "settings",
		Short:             "设置",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if cmd.Flags().Changed("aria2-enabled") {
				configs.Configs.Aria2c.Enabled = flags.Aria2c.Enabled
			}
			if cmd.Flags().Changed("aria2-retry-wait") {
				configs.Configs.Aria2c.RetryWait = flags.Aria2c.RetryWait
			}
			if cmd.Flags().Changed("aria2-split") {
				configs.Configs.Aria2c.Split = flags.Aria2c.Split
			}
			if cmd.Flags().Changed("aria2-max-connection-per-server") {
				configs.Configs.Aria2c.MaxConnectionPerServer = flags.Aria2c.MaxConnectionPerServer
			}
			if cmd.Flags().Changed("auto-accept-eula") {
				configs.Configs.AutoAcceptEULA = flags.AutoAcceptEULA
			}
			if cmd.Flags().Changed("Language") {
				switch Language {
				case "en":
					configs.Configs.Language = language.English
				case "zh":
					configs.Configs.Language = language.Chinese
				default:
					fmt.Println("尚未支持该语言...")
				}
			}
			if err := configs.Configs.Save(); err != nil {
				return err
			}
			log.Info("设置成功")
			return nil
		},
	}
	cmd.Flags().BoolVar(&flags.Aria2c.Enabled, "aria2-enabled", false, "Aria2c will be used for downloading of artifacts.")
	cmd.Flags().IntVar(&flags.Aria2c.RetryWait, "aria2-retry-wait", 0, "Number of seconds to wait between retries.")
	cmd.Flags().IntVar(&flags.Aria2c.Split, "aria2-split", 0, "Number of connections used for download.")
	cmd.Flags().IntVar(&flags.Aria2c.MaxConnectionPerServer, "aria2-max-connection-per-server", 0, "The maximum number of connections to one server for each download.")
	cmd.Flags().BoolVar(&flags.AutoAcceptEULA, "auto-accept-eula", false, "自动同意EULA")
	cmd.Flags().StringVar(&Language, "Language", "", "语言")
	return cmd
}
