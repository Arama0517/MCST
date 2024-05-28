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

// 无极镜像

package apis

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

var FastMirror = map[string]FastMirrorData{}

func InitFastMirror() {
	var err error
	url := url.URL{
		Scheme: "https",
		Host:   "download.fastmirror.net",
		Path:   "/api/v3",
	}
	resp, err := lib.Request(url, http.MethodGet, nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var data struct {
		Data []FastMirrorData `json:"data"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(data.Data); i++ {
		FastMirror[data.Data[i].Name] = data.Data[i]
	}
}

type FastMirrorData struct {
	Name        string   `json:"name"`
	Tag         string   `json:"tag"`
	Homepage    string   `json:"homepage"`
	Recommanded bool     `json:"recommanded"`
	MC_Versions []string `json:"mc_versions"`
}

type FastMirrorBuilds struct {
	Name         string `json:"name"`
	MC_Version   string `json:"mc_version"`
	Core_Version string `json:"core_version"`
	Update_Time  string `json:"update_time"`
	Sha1         string `json:"sha1"`
}

type ParsedFastMirrorBuilds map[string]FastMirrorBuilds

func GetFastMirrorBuildsDatas(ServerType string, MinecraftVersion string) (ParsedFastMirrorBuilds, error) {
	resp, err := lib.Request(url.URL{
		Scheme:   "https",
		Host:     "download.fastmirror.net",
		Path:     "/api/v3/" + ServerType + "/" + MinecraftVersion,
		RawQuery: "?offset=0&limit=25",
	}, http.MethodGet, nil)
	if err != nil {
		return ParsedFastMirrorBuilds{}, err
	}
	defer resp.Body.Close()
	var data struct {
		Data struct {
			Builds []FastMirrorBuilds `json:"builds"`
		} `json:"data"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return ParsedFastMirrorBuilds{}, err
	}
	parseDatas := ParsedFastMirrorBuilds{}
	for i := 0; i < len(data.Data.Builds); i++ {
		data := data.Data.Builds[i]
		parseDatas[data.Core_Version] = data
	}
	return parseDatas, nil
}

func DownloadFastMirrorServer(info lib.ServerInfo) (string, error) {
	path, err := lib.Download(url.URL{
		Scheme: "https",
		Host:   "download.fastmirror.net",
		Path:   "/download/" + info.ServerType + "/" + info.MinecraftVersion + "/" + info.BuildVersion,
	}, info.ServerType+"-"+info.MinecraftVersion+"-"+info.BuildVersion+".jar")
	if err != nil {
		return "", err
	}
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha1.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	hash := hasher.Sum(nil)
	FastMirrorBuildsData, err := GetFastMirrorBuildsDatas(info.ServerType, info.MinecraftVersion)
	if err != nil {
		return "", err
	}
	if fmt.Sprintf("%x", hash) != FastMirrorBuildsData[info.BuildVersion].Sha1 {
		err := os.Remove(path)
		if err != nil {
			return "", err
		}
		return "", errors.New("Sha1不匹配")
	}
	return path, nil
}
