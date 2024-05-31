/*
 * Minecraft Server Tool(MCST) is a command-line utility making Minecraft server creation quick and easy for beginners.
 * Copyright (C) 2024 Arama
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
	"os"

	"github.com/Arama-Vanarana/MCServerTool/internal/pages"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

func init() {
	lib.InitData()
	lib.InitAria2c()
}

func main() {
	app := cli.App{
		Name:    "MCST",
		Usage:   "Minecraft Server Tool",
		Version: lib.Version,
		Commands: []*cli.Command{
			&pages.Create,
			&pages.Download,
			&pages.List,
			&pages.Start,
			&pages.Config,
			&pages.Completion,
		},
		EnableBashCompletion: true,
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
