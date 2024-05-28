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

package pages

import (
	"fmt"

	"github.com/Arama-Vanarana/MCSCS-Go/apis"
	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

func DownloadPage() error {
	var info lib.ServerInfo
	var err error
serverType:
	info.ServerType, err = serverType()
	if err != nil {
		return err
	}
	if info.ServerType == "" {
		return nil
	}
minecraftVersion:
	info.MinecraftVersion, err = minecraftVersion(info.ServerType)
	if err != nil {
		return err
	}
	if info.MinecraftVersion == "" {
		goto serverType
	}
	info.BuildVersion, err = buildVersion(info.ServerType, info.MinecraftVersion)
	if err != nil {
		return err
	}
	if info.BuildVersion == "" {
		goto minecraftVersion
	}
	path, err := apis.DownloadFastMirrorServer(info)
	if err != nil {
		return err
	}
	DownloadInfo, err := lib.LoadDownloadsLists()
	if err != nil {
		return err
	}
	DownloadInfo = append(DownloadInfo, lib.DownloadInfo{
		Path: path,
		Info: info,
	})
	lib.SaveDownloadsLists(DownloadInfo)
	return nil
}

func serverType() (string, error) {
	options := []string{}
	serverTypes := []string{}
	for _, v := range apis.FastMirror {
		options = append(options, fmt.Sprintf("%s(标签: %s)", v.Name, func() string {
			tag := v.Tag
			switch tag {
			case "mod":
				return "模组"
			case "proxy":
				return "代理"
			case "bedrock":
				return "基岩"
			case "pure":
				return "纯净"
			case "vanilla":
				return "原版"
			default:
				return "未知"
			}
		}()))
		serverTypes = append(serverTypes, v.Name)
	}
	options = append(options, "返回")
	selection := lib.Select(options, "请选择一个使用的服务器类型")
	switch selection {
	case len(options):
		return "", nil
	default:
		return serverTypes[selection], nil
	}
}

func minecraftVersion(serverType string) (string, error) {
	options := apis.FastMirror[serverType].MC_Versions
	options = append(options, "返回")
	selection := lib.Select(options, "请选择一个Minecraft版本")
	switch selection {
	case len(options):
		return "", nil
	default:
		return options[selection], nil
	}
}

func buildVersion(serverType string, minecraftVersion string) (string, error) {
	FastMirror, err := apis.GetFastMirrorBuildsDatas(serverType, minecraftVersion)
	if err != nil {
		return "", err
	}
	options := []string{}
	buildVersions := []string{}
	for _, v := range FastMirror {
		options = append(options, fmt.Sprintf("%s(更新时间: %s)", v.Core_Version, v.Update_Time))
		buildVersions = append(buildVersions, v.Core_Version)
	}
	options = append(options, "返回")
	selection := lib.Select(options, "请选择一个构建版本版本")
	switch selection {
	case len(options):
		return "", nil
	default:
		return buildVersions[selection], nil
	}
}
