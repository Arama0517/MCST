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

package main

import (
	"os"

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/download"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:                   "download [url]",
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		PersistentPreRunE: func(*cobra.Command, []string) error {
			log.SetHandler(cli.Default)
			return configs.InitData()
		},
		RunE: func(_ *cobra.Command, args []string) error {
			path, err := download.NewDownloader(args[0]).Download()
			if err != nil {
				return err
			}
			log.Info(path)
			return nil
		},
	}
	if err := cmd.Execute(); err != nil {
		log.WithError(err).Fatal("下载失败")
		os.Exit(1)
	}
}
