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
	"os"
	"reflect"

	"github.com/Arama0517/MCST/internal/build"
	"github.com/Arama0517/MCST/internal/configs"
	MCSTErrors "github.com/Arama0517/MCST/internal/errors"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/Arama0517/MCST/pkg/cmd/create"
	"github.com/Arama0517/MCST/pkg/cmd/settings"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

var ExitFunc = os.Exit

func init() {
	log.SetHandler(cli.Default)
	if err := configs.InitData(); err != nil {
		log.WithError(err).Fatal("初始化配置失败")
		ExitFunc(MCSTErrors.InitConfigFail)
	}
	if err := locale.InitLocale(); err != nil {
		log.WithError(err).Error("初始化语言失败")
		ExitFunc(MCSTErrors.InitLocaleFail)
	}
}

func Execute(args []string) {
	cmd := newRootCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil && err.Error() != "interrupt" {
		log.WithError(err).Error("出现错误!")
		ExitFunc(MCSTErrors.RunFail)
	}
}

func newRootCmd() *cobra.Command {
	var debug bool
	cmd := &cobra.Command{
		Use:               "MCST",
		Short:             locale.GetLocaleMessage("root.short"),
		Long:              build.Version.ASCIIName,
		Version:           build.Version.String(),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		PersistentPreRun: func(*cobra.Command, []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}
		},
		PostRun: func(*cobra.Command, []string) {
			log.Info("运行成功")
		},
	}
	cmd.SetVersionTemplate("{{.Version}}")
	cmd.PersistentFlags().BoolVar(&debug, "debug", false, locale.GetLocaleMessage("root.flags.debug"))
	cmd.AddCommand(
		create.New(),
		newDownloadCmd(),
		newConfigCmd(),
		newStartCmd(),
		newListCmd(),
		settings.New(),
		newManCmd(),
	)
	return cmd
}

func structToMap(obj any, parentKey string, result map[string]any) {
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 构建完整的键
		tag := field.Tag.Get("yaml")
		if tag == "" {
			tag = field.Name
		}
		key := tag
		if parentKey != "" {
			key = parentKey + "." + tag
		}

		switch value.Kind() {
		case reflect.Struct:
			structToMap(value.Interface(), key, result)
		case reflect.Map:
			continue
		default:
			result[key] = value.Interface()
		}
	}
}
