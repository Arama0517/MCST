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
	"runtime"
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
	downloaderIDM     string
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
			downloaderIDM = locale.GetLocaleMessage("settings.downloader.IDM")

			aria2RetryWait = locale.GetLocaleMessage("settings.aria2.retry-wait")
			aria2Split = locale.GetLocaleMessage("settings.aria2.split")
			aria2MaxConnectionPerServer = locale.GetLocaleMessage("settings.aria2.max-connection-per-server")
		},
		RunE: func(*cobra.Command, []string) error {
			var result string
			prompt := &survey.Select{
				Message: "请选择你要设置的选项:",
				Options: []string{
					language,
					downloader,
					aria2,
					autoAcceptEULA,
				},
			}
			if err := survey.AskOne(prompt, &result); err != nil {
				return err
			}
			switch result {
			case language:
				if err := caseLanguage(); err != nil {
					return err
				}
			case downloader:
				if err := caseDownloader(); err != nil {
					return err
				}
			case aria2:
				if err := caseAria2(); err != nil {
					return err
				}
			case autoAcceptEULA:
				if err := caseAutoAcceptEULA(); err != nil {
					return err
				}
			}
			return configs.Configs.Save()
		},
	}
	return cmd
}

func caseLanguage() error {
	result := ""
	prompt := &survey.Select{
		Message: "请选择要设置的语言:",
		Options: []string{
			languageCN,
			languageEN,
		},
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return err
	}
	switch result {
	case languageCN:
		configs.Configs.Settings.Language = "zh"
	case languageEN:
		configs.Configs.Settings.Language = "en"
	}
	return nil
}

func caseDownloader() error {
	result := ""
	options := []string{
		downloaderDefault,
		downloaderAria2,
	}
	if runtime.GOOS == "windows" {
		options = append(options, downloaderIDM)
	}
	prompt := &survey.Select{
		Message: "请选择要使用的下载器:",
		Options: options,
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return err
	}
	switch result {
	case downloaderDefault:
		configs.Configs.Settings.Downloader = 0
	case downloaderAria2:
		configs.Configs.Settings.Downloader = 1
	case downloaderIDM:
		configs.Configs.Settings.Downloader = 2
	}
	return nil
}

func aria2Validator(ans any) error {
	str, ok := ans.(string)
	if !ok {
		return errors.New("无效")
	}
	_, err := strconv.Atoi(str)
	return err
}

func caseAria2() error {
	var result string
	prompt := &survey.Select{
		Message: "请选择一个Aria2配置项:",
		Options: []string{
			aria2RetryWait,
			aria2Split,
			aria2MaxConnectionPerServer,
		},
	}
	if err := survey.AskOne(prompt, &result); err != nil {
		return err
	}

	switch result {
	case aria2RetryWait:
		prompt := &survey.Input{
			Message: "请输入 '--retry-wait' 的值",
		}
		if err := survey.AskOne(prompt, &result, survey.WithValidator(aria2Validator)); err != nil {
			return err
		}
		result, err := strconv.Atoi(result)
		if err != nil {
			return err
		}
		configs.Configs.Settings.Aria2.RetryWait = result
	case aria2Split:
		prompt := &survey.Input{
			Message: "请输入 '--split' 的值",
		}
		if err := survey.AskOne(prompt, &result, survey.WithValidator(aria2Validator)); err != nil {
			return err
		}
		result, err := strconv.Atoi(result)
		if err != nil {
			return err
		}
		configs.Configs.Settings.Aria2.Split = result
	case aria2MaxConnectionPerServer:
		prompt := &survey.Input{
			Message: "请输入 '--max-connection-per-server' 的值",
		}
		if err := survey.AskOne(prompt, &result, survey.WithValidator(aria2Validator)); err != nil {
			return err
		}
		result, err := strconv.Atoi(result)
		if err != nil {
			return err
		}
		configs.Configs.Settings.Aria2.MaxConnectionPerServer = result
	}
	return nil
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
	return nil
}
