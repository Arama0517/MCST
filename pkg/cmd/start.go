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
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/apex/log"
	"github.com/go-cmd/cmd"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "start",
		Short:             locale.GetLocaleMessage("start.short"),
		Long:              locale.GetLocaleMessage("start.long"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			serverName := ""
			var serverNames []string
			for _, server := range configs.Configs.Servers {
				serverNames = append(serverNames, server.Name)
			}
			if err := survey.AskOne(&survey.Select{
				Message: "请选择一个服务器",
				Options: serverNames,
			}, &serverName); err != nil {
				return err
			}
			config := configs.Configs.Servers[serverName]

			javaCmd := cmd.NewCmd(config.Java.Path)
			javaCmd.Dir = filepath.Join(configs.ServersDir, config.Name)
			javaCmd.Args = append(javaCmd.Args,
				fmt.Sprintf("-Xms%d", config.Java.MinMemory),
				fmt.Sprintf("-Xmx%d", config.Java.MaxMemory),
				fmt.Sprintf("-Dfile.encoding=%s", config.Java.Encoding))
			javaCmd.Args = append(javaCmd.Args, config.Java.Args...)
			javaCmd.Args = append(javaCmd.Args, "-jar", "server.jar")
			javaCmd.Args = append(javaCmd.Args, config.ServerArgs...)
			stdout := make(chan string, 1)
			stderr := make(chan string, 1)
			javaCmd.Stdout = stdout
			javaCmd.Stderr = stderr
			statusChan := javaCmd.Start()

			for {
				select {
				case finalStatus := <-statusChan:
					if finalStatus.Exit != 0 {
						log.WithField("错误代码", finalStatus.Exit).Error("服务器进程未正常退出")
						ExitFunc(finalStatus.Exit)
					}
					return nil
				case str := <-stdout:
					log.Info(str)
				case str := <-stderr:
					log.Error(str)
				}
			}
		},
	}
}
