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

package api_test

import (
	"encoding/json"
	"testing"

	api "github.com/Arama-Vanarana/MCServerTool/pkg/API"
	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
)

func init() {
	if err := lib.InitAll(); err!= nil {
		panic(err)
	}
}

func TestFastMirror(t *testing.T) {
	data, err := api.GetFastMirrorDatas()
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func TestFastMirrorBuilds(t *testing.T) {
	data, err := api.GetFastMirrorDatas()
	if err != nil {
		t.Fatal(err)
	}
	MinecraftVersion := data["Mohist"].MinecraftVersions[0]
	builds, err := api.GetFastMirrorBuildsDatas("Mohist", MinecraftVersion)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.MarshalIndent(builds, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func TestFastMirrorBuildsDownload(t *testing.T) {
	data, err := api.GetFastMirrorDatas()
	if err != nil {
		t.Fatal(err)
	}
	MinecraftVersion := data["Mohist"].MinecraftVersions[0]
	builds, err := api.GetFastMirrorBuildsDatas("Mohist", MinecraftVersion)
	if err != nil {
		t.Fatal(err)
	}
	var build string
	for k := range builds {
		build = k
		break
	}
	path, err := api.DownloadFastMirrorServer("Mohist", MinecraftVersion, build)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
}
