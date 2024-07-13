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
	"fmt"

	"github.com/Arama0517/MCST/internal/dialog/files"
	"github.com/spf13/cobra"
)

func main() {
	var fileName string
	cmd := &cobra.Command{
		Use:   "choice",
		Short: "测试选择文件",
		RunE: func(*cobra.Command, []string) error {
			path, err := files.Run(fileName, func(string) bool { return true })
			if err != nil {
				return err
			}
			fmt.Println(path)
			return nil
		},
	}
	cmd.Flags().StringVarP(&fileName, "file-name", "n", "", "文件名")
	_ = cmd.MarkFlagRequired("file-name")
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
