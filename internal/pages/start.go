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
	"errors"
	"fmt"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
	"os/exec"
	"path/filepath"
)

var Start = cli.Command{
	Name:    "start",
	Aliases: []string{"s"},
	Usage:   "启动服务器",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "server",
			Aliases:  []string{"s"},
			Usage:    "要启动的服务器名称",
			Required: true,
		},
	},
	Action: func(context *cli.Context) error {
		configs, err := lib.LoadConfigs()
		if err != nil {
			return err
		}
		server, exists := configs.Servers[context.String("server")]
		if !exists {
			return errors.New("服务器不存在")
		}
		cmd := exec.Command(server.Java.Path)
		cmd.Args = append(cmd.Args, fmt.Sprintf("-Xmx%d", server.Java.Xmx), fmt.Sprintf("-Xms%d", server.Java.Xms), "-Dfile.encoding=UTF-8")
		cmd.Args = append(cmd.Args, server.Java.Args...)
		cmd.Args = append(cmd.Args, "-jar", "server.jar")
		cmd.Args = append(cmd.Args, server.ServerArgs...)
		cmd.Dir = filepath.Join(lib.ServersDir, server.Name)
		cmd.Stdout = context.App.Writer
		cmd.Stderr = context.App.ErrWriter
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	},
}
