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
	"strconv"
	"strings"

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

func newSettingsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "settings",
		Short:             locale.GetLocaleMessage("settings"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		Run: func(*cobra.Command, []string) {
			result := make(map[string]any)
			structToMap(configs.Configs.Settings, "", result)
			for k, v := range result {
				log.WithField("value", v).Info(strings.ReplaceAll(k, "_", "-"))
			}
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:  "reset",
		Long: locale.GetLocaleMessage("settings.reset"),
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			configs.Configs.Settings = configs.DefaultSettings
			return configs.Configs.Save()
		},
	})
	cmd.AddCommand(
		&cobra.Command{
			Use:   "aria2-enabled",
			Short: locale.GetLocaleMessage("settings.downloader"),
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				fmt.Println(args)
				if len(args) == 0 {
					log.WithField("value", configs.Configs.Settings.Downloader).Info("downloader")
					return nil
				}
				value, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				configs.Configs.Settings.Downloader = value
				return configs.Configs.Save()
			},
		},
		&cobra.Command{
			Use:   "aria2-retry-wait",
			Short: locale.GetLocaleMessage("settings.aria2.retry-wait"),
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if len(args) == 0 {
					log.WithField("value", configs.Configs.Settings.Aria2.RetryWait).Info("aria2.retry-wait")
					return nil
				}
				value, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				configs.Configs.Settings.Aria2.RetryWait = value
				return configs.Configs.Save()
			},
		},
		&cobra.Command{
			Use:   "aria2-split",
			Short: locale.GetLocaleMessage("settings.aria2.split"),
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if len(args) == 0 {
					log.WithField("value", configs.Configs.Settings.Aria2.Split).Info("aria2.split")
					return nil
				}
				value, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				configs.Configs.Settings.Aria2.Split = value
				return configs.Configs.Save()
			},
		},
		&cobra.Command{
			Use:   "aria2-max-connection-per-server",
			Short: locale.GetLocaleMessage("settings.aria2.max-connection-per-server"),
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if len(args) == 0 {
					log.WithField("value", configs.Configs.Settings.Aria2.MaxConnectionPerServer).Info("aria2.max-connection-per-server")
					return nil
				}
				value, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				configs.Configs.Settings.Aria2.MaxConnectionPerServer = value
				return configs.Configs.Save()
			},
		},
		newSettingsAria2OptionsCmd(),
		&cobra.Command{
			Use:   "auto-accept-eula",
			Short: locale.GetLocaleMessage("settings.auto-accept-eula"),
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if len(args) == 0 {
					log.WithField("value", configs.Configs.Settings.AutoAcceptEULA).Info("auto-accept-eula")
					return nil
				}
				value, err := strconv.ParseBool(args[0])
				if err != nil {
					return err
				}
				configs.Configs.Settings.AutoAcceptEULA = value
				return configs.Configs.Save()
			},
		},
		&cobra.Command{
			Use:   "language",
			Short: locale.GetLocaleMessage("settings.language"),
			Args:  cobra.MaximumNArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				if len(args) == 0 {
					log.WithField("value", configs.Configs.Settings.Language).Info("language")
				}
				switch args[0] {
				case language.English.String():
					configs.Configs.Settings.Language = language.English.String()
				case language.Chinese.String():
					configs.Configs.Settings.Language = language.Chinese.String()
				default:
					log.Warn(locale.GetLocaleMessage("errors.language"))
				}
				return configs.Configs.Save()
			},
		})
	// cmd.Flags().BoolVar(&flags.Aria2.Enabled, "aria2-enabled", false, locale.GetLocaleMessage("settings.aria2-enabled"))
	// cmd.Flags().IntVar(&flags.Aria2.RetryWait, "aria2-retry-wait", 0, locale.GetLocaleMessage("settings.aria2-retry-wait"))
	// cmd.Flags().IntVar(&flags.Aria2.Split, "aria2-split", 0, locale.GetLocaleMessage("settings.aria2-split"))
	// cmd.Flags().IntVar(&flags.Aria2.MaxConnectionPerServer, "aria2-max-connection-per-server", 0, locale.GetLocaleMessage("settings.aria2-max-connection-per-server"))
	// cmd.Flags().BoolVar(&flags.AutoAcceptEULA, "auto-accept-eula", false, locale.GetLocaleMessage("settings.auto-accept-eula"))
	// cmd.Flags().StringVar(&Language, "language", "", locale.GetLocaleMessage("settings.language"))
	return cmd
}

func newSettingsAria2OptionsCmd() *cobra.Command {
	var options []string
	cmd := &cobra.Command{
		Use:   "aria2-options",
		Short: "",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("value") {
				log.WithField("value", configs.Configs.Settings.Aria2.Options).Info("aria2.options")
			}
			configs.Configs.Settings.Aria2.Options = options
			return configs.Configs.Save()
		},
	}
	cmd.Flags().StringArrayVar(&options, "value", nil, "")
	return cmd
}
