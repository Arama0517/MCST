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
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Arama0517/MCST/internal/configs"
	MCSTErrors "github.com/Arama0517/MCST/internal/errors"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/spf13/cobra"
)

func newStartCmd() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:               "start",
		Short:             locale.GetLocaleMessage("start.short"),
		Long:              locale.GetLocaleMessage("start.long"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			config, exists := configs.Configs.Servers[name]
			if !exists {
				return MCSTErrors.ErrServerNotFound
			}
			cmd := exec.Command(config.Java.Path)
			cmd.Dir = filepath.Join(configs.ServersDir, config.Name)
			cmd.Args = append(cmd.Args,
				fmt.Sprintf("-Xms%d", config.Java.Xms),
				fmt.Sprintf("-Xmx%d", config.Java.Xmx),
				fmt.Sprintf("-Dfile.encoding=%s", config.Java.Encoding))
			cmd.Args = append(cmd.Args, config.Java.Args...)
			cmd.Args = append(cmd.Args, "-jar", "server.jar")
			cmd.Args = append(cmd.Args, config.ServerArgs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		},
	}
	cmd.Flags().StringVarP(&name, "name", "n", "", locale.GetLocaleMessage("start.flags.name"))
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
