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
	_ "embed"
	"os"

	"github.com/Arama-Vanarana/MCServerTool/cmd"
)

var version = "dev" // 开发时使用的版本号, 发布时使用 '-X main.version=X.X.X' 参数构建

func main() {
	err := cmd.Execute(os.Exit, os.Args[1:], version)
	if err != nil {
		return
	}
}
