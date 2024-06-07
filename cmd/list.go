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

	"github.com/Arama-Vanarana/MCServerTool/internal/lib"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "列出服务器",
		Long:  "列出所有服务器的名称",
		RunE: func(_ *cobra.Command, _ []string) error {
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			if _, err := fmt.Fprintln(os.Stderr, "已创建的服务器:"); err != nil {
				return err
			}
			for _, config := range configs.Servers {
				fmt.Println(config.Name)
			}
			return nil
		},
	}
	return cmd
}
