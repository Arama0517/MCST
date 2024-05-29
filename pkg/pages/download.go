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

	"github.com/Arama-Vanarana/MCSCS-Go/pkg/apis"
	"github.com/Arama-Vanarana/MCSCS-Go/pkg/lib"
)

func DownloadPage() error {
	var info lib.ServerInfo
	var err error
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
serverType:
	info.ServerType = serverType()
	if info.ServerType == "" {
		return nil
	}
minecraftVersion:
	info.MinecraftVersion = minecraftVersion(info.ServerType)
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
	configs.Downloads = append(configs.Downloads, lib.DownloadInfo{
		Path: path,
		Info: info,
	})
	configs.Save()
	return nil
}

func serverType() string {
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
	selection := lib.Select("请选择一个使用的服务器类型", options)
	switch selection {
	case len(options):
		return ""
	default:
		return serverTypes[selection]
	}
}

func minecraftVersion(serverType string) string {
	options := apis.FastMirror[serverType].MC_Versions
	options = append(options, "返回")
	selection := lib.Select("请选择一个Minecraft版本", options)
	switch selection {
	case len(options):
		return ""
	default:
		return options[selection]
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
	selection := lib.Select("请选择一个构建版本版本", options)
	switch selection {
	case len(options):
		return "", nil
	default:
		return buildVersions[selection], nil
	}
}
