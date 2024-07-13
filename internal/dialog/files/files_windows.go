//go:build windows

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

package files

import (
	"github.com/apex/log"
	"github.com/lxn/win"
	"github.com/ncruces/zenity"
)

// Run 使用资源管理器原生API
func Run(fileName string, checkFunc func(path string) bool) (string, error) {
	for {
		path, err := zenity.SelectFile(
			zenity.Title("请选择文件"),
			zenity.Attach(win.GetForegroundWindow()),
			zenity.FileFilters{
				zenity.FileFilter{
					Name:     "需要的文件",
					Patterns: []string{fileName},
					CaseFold: false,
				},
				zenity.FileFilter{
					Name:     "任意文件",
					Patterns: []string{"*"},
					CaseFold: false,
				},
			},
		)
		if err != nil {
			return "", err
		}
		if checkFunc(path) {
			return path, nil
		}
		log.Warnf("选择的文件不正确或无效, 请选择(指向)名为 '%s' 的文件", fileName)
	}
}
