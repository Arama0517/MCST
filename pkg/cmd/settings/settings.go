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
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/spf13/cobra"
)

var (
	language       string
	downloader     string
	aria2          string
	autoAcceptEULA string
)

var (
	languageCN string
	languageEN string
)

var (
	downloaderDefault string
	downloaderAria2   string
)

var (
	aria2RetryWait              string
	aria2Split                  string
	aria2MaxConnectionPerServer string
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "settings",
		Short:         locale.GetLocaleMessage("settings"),
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		PersistentPreRun: func(*cobra.Command, []string) {
			language = locale.GetLocaleMessage("settings.language")
			downloader = locale.GetLocaleMessage("settings.downloader")
			aria2 = locale.GetLocaleMessage("settings.aria2")
			autoAcceptEULA = locale.GetLocaleMessage("settings.auto-accept-eula")

			languageCN = "简体中文"
			languageEN = "English"

			downloaderDefault = locale.GetLocaleMessage("settings.downloader.default")
			downloaderAria2 = locale.GetLocaleMessage("settings.downloader.aria2")

			aria2RetryWait = locale.GetLocaleMessage("settings.aria2.retry-wait")
			aria2Split = locale.GetLocaleMessage("settings.aria2.split")
			aria2MaxConnectionPerServer = locale.GetLocaleMessage("settings.aria2.max-connection-per-server")
		},
		RunE: func(*cobra.Command, []string) error {
			result := ""
			if err := survey.AskOne(&survey.Select{
				Message: "请选择你要设置的选项:",
				Options: []string{
					language,
					downloader,
					aria2,
					autoAcceptEULA,
				},
			}, &result); err != nil {
				return err
			}
			switch result {
			case language:
				return caseLanguage()
			case downloader:
				return caseDownloader()
			case aria2:
				return caseAria2()
			case autoAcceptEULA:
				return caseAutoAcceptEULA()
			}
			return nil
		},
	}
	return cmd
}

func caseLanguage() error {
	result := ""
	if err := survey.AskOne(&survey.Select{
		Message: "请选择你要使用的语言",
		Options: []string{
			languageCN,
			languageEN,
		},
	}, &result); err != nil {
		return err
	}
	switch result {
	case languageCN:
		configs.Configs.Settings.Language = "zh"
	case languageEN:
		configs.Configs.Settings.Language = "en"
	}
	return configs.Configs.Save()
}

func caseDownloader() error {
	result := ""
	prompt := &survey.Select{
		Message: "请选择要使用的下载器",
		Options: []string{
			downloaderDefault,
			downloaderAria2,
		},
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return err
	}
	switch result {
	case downloaderDefault:
		configs.Configs.Settings.Downloader = 0
	case downloaderAria2:
		configs.Configs.Settings.Downloader = 1
	}
	return configs.Configs.Save()
}

func aria2Validator(ans any) error {
	str, ok := ans.(string)
	if !ok {
		return errors.New("无效")
	}
	intAns, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	if intAns <= 0 {
		return fmt.Errorf("值不可小于等于0")
	}
	return nil
}

func caseAria2() error {
	result := ""
	if err := survey.AskOne(&survey.Select{
		Message: "请选择一个Aria2配置项",
		Options: []string{
			aria2RetryWait,
			aria2Split,
			aria2MaxConnectionPerServer,
		},
	}, &result); err != nil {
		return err
	}
	switch result {
	case aria2RetryWait:
		if err := survey.AskOne(&survey.Input{Message: "请输入值:"}, &result, survey.WithValidator(aria2Validator)); err != nil {
			return err
		}
		result, err := strconv.Atoi(result)
		if err != nil {
			return err
		}
		configs.Configs.Settings.Aria2.RetryWait = result
	case aria2Split:
		if err := survey.AskOne(&survey.Input{Message: "请输入值:"}, &result, survey.WithValidator(aria2Validator)); err != nil {
			return err
		}
		result, err := strconv.Atoi(result)
		if err != nil {
			return err
		}
		configs.Configs.Settings.Aria2.Split = result
	case aria2MaxConnectionPerServer:
		if err := survey.AskOne(&survey.Input{Message: "请输入值:"}, &result, survey.WithValidator(func(ans any) error {
			if err := aria2Validator(ans); err != nil {
				return err
			}
			intAns, err := strconv.Atoi(ans.(string))
			if err != nil {
				return err
			}
			if intAns > 16 {
				return fmt.Errorf("此参数最大值为16")
			}
			return nil
		})); err != nil {
			return err
		}
		result, err := strconv.Atoi(result)
		if err != nil {
			return err
		}
		configs.Configs.Settings.Aria2.MaxConnectionPerServer = result
	}
	return configs.Configs.Save()
}

func caseAutoAcceptEULA() error {
	result := false
	prompt := &survey.Confirm{
		Message: "是否启用自动同意EULA?",
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return err
	}
	configs.Configs.Settings.AutoAcceptEULA = result
	return configs.Configs.Save()
}
