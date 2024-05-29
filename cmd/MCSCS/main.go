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
	api "github.com/Arama-Vanarana/MCSCS-Go/pkg/API"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/pages"
	"github.com/rs/zerolog/log"
)

func init() {
	lib.Init()
	api.Init()
}

func main() {
	options := []string{"创建服务器", "下载核心", "退出"}
	for {
		selected := lib.Select("请选择一个选项", options)
		var err error
		switch selected {
		case 0:
			err = pages.CreatePage()
		case 1:
			err = pages.DownloadPage()
		default:
			return
		}
		if err != nil {
			log.Error().Err(err).Msg("运行错误")
		}
	}

}