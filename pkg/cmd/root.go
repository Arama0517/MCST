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
	"github.com/Arama0517/MCST/internal/build"
	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

func Execute(args []string, exit func(code int)) {
	log.SetHandler(cli.Default)
	if err := configs.InitData(); err != nil {
		log.WithError(err).Fatal("初始化配置失败")
		exit(InitConfigFail)
	}
	if err := locale.InitLocale(); err != nil {
		log.WithError(err).Error("初始化语言失败")
		exit(InitLocaleFail)
	}
	cmd := newRootCmd()
	cmd.SetArgs(args)
	if err := cmd.Execute(); err != nil {
		log.WithError(err).Error("出现错误!")
		exit(RunFail)
	}
	exit(OK)
}

func newRootCmd() *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:               "MCST",
		Short:             locale.GetLocaleMessage("root.short"),
		Long:              build.Version.ASCIIName,
		Version:           build.Version.String(),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			if verbose {
				log.SetLevel(log.DebugLevel)
			}
			return nil
		},
		PostRun: func(*cobra.Command, []string) {
			log.Info("运行成功")
		},
	}
	cmd.SetVersionTemplate("{{.Version}}")
	cmd.PersistentFlags().BoolVar(&verbose, "debug", false, locale.GetLocaleMessage("root.flags.debug"))
	cmd.AddCommand(
		newCreateCmd(),
		newDownloadCmd(),
		newConfigCmd(),
		newStartCmd(),
		newListCmd(),
		newSettingsCmd(),
		newManCmd(),
	)
	return cmd
}
