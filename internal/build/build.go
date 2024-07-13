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

package build

import (
	_ "embed"
	"net/url"

	goversion "github.com/caarlos0/go-version"
)

var (
	version   = ""
	commit    = ""
	treeState = ""
	date      = ""
	builtBy   = ""
)

//go:embed art.txt
var asciiArt string

var Version goversion.Info

func init() {
	Version = goversion.GetVersionInfo(func(i *goversion.Info) {
		i.Name = "Minecraft Server Tool"
		i.Description = "A command-line utility making Minecraft server creation quick and easy for beginners."
		i.URL = (&url.URL{
			Scheme: "https",
			Host:   "github.com",
			Path:   "/Arama0517/MCST",
		}).String()
		i.ASCIIName = asciiArt
		if commit != "" {
			i.GitCommit = commit
		}
		if treeState != "" {
			i.GitTreeState = treeState
		}
		if date != "" {
			i.BuildDate = date
		}
		if version != "" {
			i.GitVersion = version
		}
		if builtBy != "" {
			i.BuiltBy = builtBy
		}
	})
}
