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
	"github.com/Arama-Vanarana/MCServerTool/cmd/download"
	"github.com/Arama-Vanarana/MCServerTool/internal/lib"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

func Execute(exit func(int), args []string, version string) error {
	cmd := newRootCmd(version)
	cmd.SetArgs(args)
	log.WithField("args", args).Debug("参数")
	if err := cmd.Execute(); err != nil {
		log.WithError(err).Error("错误")
		exit(1)
	}
	return nil
}

func newRootCmd(version string) *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:           "MCST",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			log.SetHandler(cli.Default)
			if verbose {
				log.SetLevel(log.DebugLevel)
				log.Debug("调试模式开启")
			}
			return lib.InitAll(version)
		},
	}
	cmd.SetVersionTemplate("{{.Version}}")
	cmd.PersistentFlags().BoolVar(&verbose, "debug", false, "调试模式(更多的日志)")
	cmd.AddCommand(newCreateCmd(), download.NewDownloadCmd(), newConfigCmd(), newStartCmd(), newListCmd(), newManCmd())
	return cmd
}
