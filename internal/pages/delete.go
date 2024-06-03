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

package pages

import (
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

var Delete = cli.Command{
	Name: "删除服务器",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "server",
			Aliases:  []string{"s"},
			Usage:    "要删除的服务器名称",
			Required: true,
		},
	},
	Action: func(context *cli.Context) error {
		server := context.String("server")
		configs, err := lib.LoadConfigs()
		if err != nil {
			return err
		}
		switch c, err := confirm("确认删除此服务器?", context); {
		case err != nil:
			return err
		case c:
			delete(configs.Servers, server)
			err = configs.Save()
			if err != nil {
				return err
			}
		case !c:
			break
		}
		return nil
	},
}
