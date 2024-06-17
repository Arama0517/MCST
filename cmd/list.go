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

	"github.com/Arama0517/MCST/pkg/lib"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:               "list",
		Short:             "列出服务器",
		Long:              "列出所有服务器的名称",
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(_ *cobra.Command, _ []string) error {
			if _, err := fmt.Fprintln(os.Stderr, "已创建的服务器:"); err != nil {
				return err
			}
			for _, config := range lib.Configs.Servers {
				log.WithFields(log.Fields{
					"服务器编码":     config.Java.Encoding,
					"JVM初始堆内存":  config.Java.Xms,
					"JVM最大堆内存":  config.Java.Xmx,
					"Java虚拟机参数": config.Java.Args,
					"服务器参数":     config.ServerArgs,
				})
			}
			return nil
		},
	}
}
