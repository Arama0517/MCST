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

package apis

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

type FastMirror struct {
	Name        string   `json:"name"`
	Tag         string   `json:"tag"`
	Homepage    string   `json:"homepage"`
	Recommanded bool     `json:"recommanded"`
	MC_Versions []string `json:"mc_versions"`
}

type ParsedFastMirror map[string]FastMirror

func GetFastMirrorDatas() (ParsedFastMirror, error) {
	url := url.URL{
		Scheme: "https",
		Host:   "download.fastmirror.net",
		Path:   "/api/v3",
	}
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return ParsedFastMirror{}, err
	}
	req.Header.Set("User-Agent", "MCSCS-Golang/"+lib.VERSION)
	resp, err := client.Do(req)
	if err != nil {
		return ParsedFastMirror{}, err
	}
	defer resp.Body.Close()
	var data struct {
		Data []FastMirror `json:"data"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ParsedFastMirror{}, err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return ParsedFastMirror{}, err
	}
	parsedData := ParsedFastMirror{}
	for i := 0; i < len(data.Data); i++ {
		data := data.Data[i]
		parsedData[data.Name] = FastMirror{
			Name:        data.Name,
			Tag:         data.Tag,
			Homepage:    data.Homepage,
			Recommanded: data.Recommanded,
			MC_Versions: data.MC_Versions,
		}
	}
	return parsedData, nil
}

type FastMirrorBuilds struct {
	Name         string `json:"name"`
	MC_Version   string `json:"mc_version"`
	Core_Version string `json:"core_version"`
	Update_Time  string `json:"update_time"`
	Sha1         string `json:"sha1"`
}

type ParsedFastMirrorBuilds map[string]FastMirrorBuilds

func GetFastMirrorBuildsDatas(server_type string, minecraft_version string) (ParsedFastMirrorBuilds, error) {
	url := url.URL{
		Scheme:   "https",
		Host:     "download.fastmirror.net",
		Path:     "/api/v3/" + server_type + "/" + minecraft_version,
		RawQuery: "?offset=0&limit=25",
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return ParsedFastMirrorBuilds{}, err
	}
	req.Header.Set("User-Agent", "MCSCS-Golang/"+lib.VERSION)
	resp, err := client.Do(req)
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
		parseDatas[data.Core_Version] = FastMirrorBuilds{
			Name:         data.Name,
			MC_Version:   data.MC_Version,
			Core_Version: data.Core_Version,
			Update_Time:  data.Update_Time,
			Sha1:         data.Sha1,
		}
	}
	return parseDatas, nil
}

func DownloadFastMirrorServer(server_type string, minecraft_version string, build_version string) (string, error) {
	url := url.URL{
		Scheme: "https",
		Host:   "download.fastmirror.net",
		Path:   "/download/" + server_type + "/" + minecraft_version + "/" + build_version,
	}
	path, err := lib.Download(url.String(), server_type+"-"+minecraft_version+"-"+build_version+".jar")
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
	FastMirrorBuildsData, err := GetFastMirrorBuildsDatas(server_type, minecraft_version)
	if err != nil {
		return "", err
	}
	if string(hash) != FastMirrorBuildsData[build_version].Sha1 {
		return "", errors.New("Sha1不匹配")
	}
	return path, nil
}
