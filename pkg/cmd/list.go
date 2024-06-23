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

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "list",
		Short:             locale.GetLocaleMessage("list"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(*cobra.Command, []string) error {
			if _, err := fmt.Fprintln(os.Stderr, "已创建的服务器:"); err != nil {
				return err
			}
			for _, config := range configs.Configs.Servers {
				log.WithFields(log.Fields{
					locale.GetLocaleMessage("create.flags.xms"):         config.Java.Xms,
					locale.GetLocaleMessage("create.flags.xmx"):         config.Java.Xmx,
					locale.GetLocaleMessage("create.flags.encoding"):    config.Java.Encoding,
					locale.GetLocaleMessage("create.flags.java"):        config.Java.Path,
					locale.GetLocaleMessage("create.flags.jvm_args"):    config.Java.Args,
					locale.GetLocaleMessage("create.flags.server_args"): config.ServerArgs,
				}).Info(config.Name)
			}
			return nil
		},
	}
}
