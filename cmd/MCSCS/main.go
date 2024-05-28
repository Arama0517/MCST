/*
 * MCSCS can be used to easily create, launch, and configure a Minecraft server.
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
	"flag"
	"fmt"

	"github.com/Arama-Vanarana/MCSCS-Go/pkg/apis"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/pages"
)

func init() {
	lib.Init()
	apis.Init()
}

func main() {
	version := flag.Bool("version", false, "显示程序版本")
	flag.Parse()
	if *version {
		fmt.Println(lib.VERSION)
		return
	}
	options := []string{"创建服务器", "下载核心", "退出"}
	for {
		selected := lib.Select(options, "请选择一个选项 ")
		var err error
		switch selected {
		case 0:
			err = pages.CreatePage()
		case 1:
			err = pages.DownloadPage()
		default:
			return
		}
		lib.ClearScreen()
		if err != nil {
			lib.Logger.WithError(err).Error("运行失败")
		}
	}

}
