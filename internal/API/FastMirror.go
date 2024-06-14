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

// 无极镜像

package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Arama0517/MCST/pkg/lib"
)

func GetFastMirrorData() (map[string]FastMirrorData, error) {
	var err error
	resp, err := lib.Request(url.URL{
		Scheme: "https",
		Host:   "download.fastmirror.net",
		Path:   "/api/v3",
	}, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	var data struct {
		Data []FastMirrorData `json:"data"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	result := map[string]FastMirrorData{}
	for i := 0; i < len(data.Data); i++ {
		result[data.Data[i].Name] = data.Data[i]
	}
	return result, nil
}

func GetFastMirrorBuildsData(core string, minecraftVersion string) (map[string]FastMirrorBuilds, error) {
	resp, err := lib.Request(url.URL{
		Scheme:   "https",
		Host:     "download.fastmirror.net",
		Path:     "/api/v3/" + core + "/" + minecraftVersion,
		RawQuery: "offset=0&limit=25",
	}, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	var data struct {
		Data struct {
			Builds []FastMirrorBuilds `json:"builds"`
		} `json:"data"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	parseData := map[string]FastMirrorBuilds{}
	for i := 0; i < len(data.Data.Builds); i++ {
		data := data.Data.Builds[i]
		parseData[data.CoreVersion] = data
	}
	return parseData, nil
}

func GetFastMirrorDownloader(core, minecraftVersion, buildVersion string) *lib.Downloader {
	return lib.NewDownloader(url.URL{
		Scheme: "https",
		Host:   "download.fastmirror.net",
		Path:   fmt.Sprintf("/download/%s/%s/%s", core, minecraftVersion, buildVersion),
	})
}

type FastMirrorData struct {
	Name              string   `json:"name"`
	Tag               string   `json:"tag"`
	Homepage          string   `json:"homepage"`
	Recommend         bool     `json:"recommend"`
	MinecraftVersions []string `json:"mc_versions"`
}

type FastMirrorBuilds struct {
	Name             string `json:"name"`
	MinecraftVersion string `json:"mc_version"`
	CoreVersion      string `json:"core_version"`
	UpdateTime       string `json:"update_time"`
	Sha1             string `json:"sha1"`
}
