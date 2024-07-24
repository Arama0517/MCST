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

package settings

import (
	"errors"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "settings",
		Short:         locale.GetLocaleMessage("settings.short"),
		Long:          locale.GetLocaleMessage("settings.long"),
		SilenceErrors: true,
		SilenceUsage:  true,
		Args:          cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			result := ""
			if err := survey.AskOne(&survey.Select{
				Message: "请选择一个要更改的设置",
				Options: []string{
					"Language",
					"Aria2",
					"Auto accept EULA",
				},
				Description: func(value string, _ int) string {
					switch value {
					case "Language":
						return locale.GetLocaleMessage("settings.language")
					case "Aria2":
						return locale.GetLocaleMessage("settings.aria2")
					case "Auto accept EULA":
						return locale.GetLocaleMessage("settings.auto_accept_eula")
					default:
						return ""
					}
				},
			}, &result); err != nil {
				return err
			}
			switch result {
			case "Language":
				return caseLanguage()
			case "Aria2":
				return caseAria2()
			case "Auto accept EULA":
				return caseAutoAcceptEULA()
			}
			return nil
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: locale.GetLocaleMessage("settings.reset.short"),
		RunE: func(*cobra.Command, []string) error {
			configs.Configs.Settings = configs.DefaultSettings
			return configs.Configs.Save()
		},
	})
	return cmd
}

var (
	Chinese = language.Chinese.String()
	English = language.English.String()
)

func caseLanguage() error {
	var result string
	if err := survey.AskOne(&survey.Select{
		Message: "请选择一种语言",
		Options: []string{
			Chinese,
			English,
		},
		Description: func(value string, _ int) string {
			switch value {
			case Chinese:
				return "简体中文"
			case English:
				return "English"
			default:
				return ""
			}
		},
	}, &result); err != nil {
		return err
	}
	switch result {
	case Chinese:
		configs.Configs.Settings.Language = Chinese
	case English:
		configs.Configs.Settings.Language = English
	}
	return configs.Configs.Save()
}

func aria2IntValidator(ans any) error {
	result, ok := ans.(string)
	if !ok {
		return errors.New("invalid answer type")
	}
	value, err := strconv.Atoi(result)
	if err != nil {
		return err
	}
	if value <= 0 {
		return errors.New("invalid answer value")
	}
	return nil
}

func caseAria2() error {
	var result string
	if err := survey.AskOne(&survey.Select{
		Message: "请选择一个Aria2的配置项",
		Options: []string{
			"enabled",
			"retry-wait",
			"split",
			"max-connection-per-server",
		},
		Description: func(value string, _ int) string {
			switch value {
			case "enabled":
				return locale.GetLocaleMessage("settings.aria2.enabled")
			case "retry-wait":
				return locale.GetLocaleMessage("settings.aria2.retry_wait")
			case "split":
				return locale.GetLocaleMessage("settings.aria2.split")
			case "max-connection-per-server":
				return locale.GetLocaleMessage("settings.aria2.max_connection_per_server")
			default:
				return ""
			}
		},
	}, &result); err != nil {
		return err
	}
	switch result {
	case "enabled":
		enabled := false
		if err := survey.AskOne(&survey.Confirm{Message: "是否启用Aria2"}, &enabled); err != nil {
			return err
		}
		configs.Configs.Settings.Aria2.Enable = enabled
	case "retry-wait":
		if err := survey.AskOne(
			&survey.Input{Message: "请输入参数的值"}, &result,
			survey.WithValidator(aria2IntValidator)); err != nil {
			return err
		}
		retryWait, _ := strconv.Atoi(result)
		configs.Configs.Settings.Aria2.RetryWait = retryWait
	case "split":
		if err := survey.AskOne(
			&survey.Input{Message: "请输入参数的值"}, &result,
			survey.WithValidator(aria2IntValidator)); err != nil {
			return err
		}
		split, _ := strconv.Atoi(result)
		configs.Configs.Settings.Aria2.Split = split
	case "max-connection-per-server":
		if err := survey.AskOne(
			&survey.Input{Message: "请输入参数的值"}, &result,
			survey.WithValidator(aria2IntValidator)); err != nil {
			return err
		}
		maxConnectionPerServer, _ := strconv.Atoi(result)
		configs.Configs.Settings.Aria2.MaxConnectionPerServer = maxConnectionPerServer
	}
	return configs.Configs.Save()
}

func caseAutoAcceptEULA() error {
	result := false
	if err := survey.AskOne(&survey.Confirm{Message: "是否启用自动同意EULA?"}, &result); err != nil {
		return err
	}
	configs.Configs.Settings.AutoAcceptEULA = result
	return configs.Configs.Save()
}
