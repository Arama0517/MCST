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
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)
	if err := cmd.cmd.Execute(); err != nil {
		log.WithError(err).Error("错误")
		cmd.exit(1)
	}
}

func Execute(exit func(int), args []string, version string) error {
	root, err := newRootCmd(exit, version)
	if err != nil {
		return err
	}
	root.Execute(args)
	return nil
}

type rootCmd struct {
	cmd     *cobra.Command
	exit    func(int)
	verbose bool
}

func newRootCmd(exit func(int), version string) (*rootCmd, error) {
	root := &rootCmd{
		exit: exit,
	}
	cmd := &cobra.Command{
		Use:               "MCST",
		Short:             "Minecraft Server Tool",
		Long:              "a command-line utility making Minecraft server creation quick and easy for beginners.",
		Version:           version,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		PersistentPreRun: func(*cobra.Command, []string) {
			if root.verbose {
				log.SetLevel(log.DebugLevel)
				log.Debug("调试模式开启")
			}
			if err := lib.InitAll(); err != nil {
				log.WithError(err).Error("初始化失败")
				root.exit(1)
			}
		},
	}
	cmd.SetVersionTemplate("{{.Version}}")
	cmd.PersistentFlags().BoolVar(&root.verbose, "debug", false, "调试模式(更多的日志)")
	createCmd, err := newCreateCmd()
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(createCmd, newConfigCmd())
	root.cmd = cmd
	return root, nil
}
