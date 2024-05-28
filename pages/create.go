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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
	"github.com/shirou/gopsutil/mem"
)

var JVM_min_ram uint64 = 1048576
var RealMem uint64
var TempConfigsPath string

func InitCreatePage() {
	RealMemInfo, err := mem.VirtualMemory()
	if err != nil {
		panic(err)
	}
	RealMem = RealMemInfo.Total
	TempConfigsPath = filepath.Join(lib.ConfigsDir, "config_temp.json")
}

func CreatePage() error {
	configs, err := lib.LoadServerConfigs()
	if err != nil {
		return err
	}
	config := lib.ServerConfig{
		// 1073741824 = 1GiB
		Ram: lib.Ram{
			XMX: 1073741824,
			XMS: 1073741824,
		},
		Encoding:   "UTF-8",
		JVMArgs:    []string{"-Dlog4j2.formatMsgNoLookups=true"},
		ServerArgs: []string{"--nogui"},
	}
	if _, err := os.Stat(TempConfigsPath); err == nil && lib.Confirm("检测到存在上次已暂存的配置, 是否还原?") {
		file, err := os.ReadFile(TempConfigsPath)
		if err != nil {
			return err
		}
		err = json.Unmarshal(file, &config)
		if err != nil {
			return err
		}
	}
main:
	for {
		lib.ClearScreen()
		options := []string{
			"服务器名称:" + func() string {
				if config.Name == "" {
					return "未设置"
				}
				return config.Name
			}(),
			"XMS(Java虚拟机初始堆内存): " + func() string {
				if config.Ram.XMS == 1073741824 {
					return fmt.Sprintf("默认配置: %dB", JVM_min_ram)
				}
				return fmt.Sprintf("%dB", config.Ram.XMS)
			}(),
			"XMX(Java虚拟机最大堆内存): " + func() string {
				if config.Ram.XMX == 1073741824 {
					return fmt.Sprintf("默认配置: %dB", JVM_min_ram)
				}
				return fmt.Sprintf("%dB", config.Ram.XMX)
			}(),
			"编码: " + func() string {
				if config.Encoding == "" {
					return "默认配置: " + config.Encoding
				}
				return config.Encoding
			}(),
			"Java虚拟机参数: " + func() string {
				if reflect.DeepEqual(config.JVMArgs, []string{"-Dlog4j2.formatMsgNoLookups=true"}) {
					return "默认配置: " + strings.Join(config.JVMArgs, " ")
				}
				return strings.Join(config.JVMArgs, " ")
			}(),
			"Java: " + func() string {
				if reflect.DeepEqual(config.Java, lib.JavaInfo{}) {
					return "未设置"
				}
				return config.Java.Version + " (" + config.Java.Path + ")"
			}(),
			"服务器参数: " + func() string {
				if reflect.DeepEqual(config.ServerArgs, []string{"--nogui"}) {
					return fmt.Sprintf("默认配置: %s", strings.Join(config.ServerArgs, " "))
				}
				return strings.Join(config.ServerArgs, " ")
			}(),
			"核心: " + func() string {
				if reflect.DeepEqual(config.Info, lib.ServerInfo{}) {
					return "未设置"
				}
				return fmt.Sprintf("%s-%s-%s.jar", config.Info.ServerType, config.Info.MinecraftVersion, config.Info.BuildVersion)
			}(),
			"返回并暂存",
			"完成并保存",
		}
		selection := lib.Select(options, "请选择一个选项")
		switch selection {
		case 0:
			config.Name = name(configs)
		case 1:
			config.Ram.XMS = jvmArgsXMS(config.Ram.XMX)
		case 2:
			config.Ram.XMX = jvmArgsXMX(config.Ram.XMS)
		case 3:
			config.Encoding = encoding()
		case 4:
			config.JVMArgs = jvmArgs(config.JVMArgs)
		case 5:
			config.Java = java()
		case 6:
			config.ServerArgs = serverArgs(config.ServerArgs)
		case 7:
			DownloadsInfo, err := lib.LoadDownloadsLists()
			if err != nil {
				return err
			}
			options := []string{}
			for _, v := range DownloadsInfo {
				options = append(options, filepath.Base(v.Path))
			}
			options = append(options, "返回")
			selection := lib.Select(options, "请选择一个服务器核心")
			if selection == len(options)-1 {
				continue
			}
			config.Info = DownloadsInfo[selection].Info
		case len(options) - 2:
			jsonConfigs, err := json.MarshalIndent(configs, "", "  ")
			if err != nil {
				return err
			}
			err = os.WriteFile(filepath.Join(lib.ConfigsDir, "config_temp.json"), jsonConfigs, 0644)
			if err != nil {
				return err
			}
			return nil
		case len(options) - 1:
			v := reflect.ValueOf(config)
			for i := 0; i < v.NumField(); i++ {
				field := v.Field(i)
				if reflect.DeepEqual(field.Interface(), reflect.Zero(reflect.TypeOf(field.Interface())).Interface()) {
					fmt.Println("你还有一些配置没有设置!")
					continue main
				}
			}
			configs[config.Name] = config
			err = lib.SaveServerConfigs(configs)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

func name(configs map[string]lib.ServerConfig) string {
main:
	for {
		inputName := lib.Input("请输入此服务器的名称: ")
		for name := range configs {
			if inputName == name {
				fmt.Println("已存在此名称的服务器, 请重新输入")
				continue main
			}
		}
		return inputName
	}
}

func ToBytes(byteStr string) (uint64, error) {
	var numPart, unitPart string

	// 分离数字部分和单位部分
	for _, char := range byteStr {
		if char >= '0' && char <= '9' {
			numPart += string(char)
		} else if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
			unitPart += strings.ToUpper(string(char))
		} else {
			return 0, nil // 如果有除了数字和单位以外的字符, 则返回0
		}
	}
	unitJSON := map[string]uint64{
		"T":     1024 * 1024 * 1024 * 1024,
		"TB":    1000 * 1000 * 1000 * 1000,
		"TIB":   1024 * 1024 * 1024 * 1024,
		"G":     1024 * 1024 * 1024,
		"GB":    1000 * 1000 * 1000,
		"GIB":   1024 * 1024 * 1024,
		"M":     1024 * 1024,
		"MB":    1000 * 1000,
		"MIB":   1024 * 1024,
		"K":     1024,
		"KB":    1000,
		"KIB":   1024,
		"B":     1,
		"BYTES": 1,
		"":      1,
	}

	unitMultiplier, ok := unitJSON[unitPart]
	if !ok {
		return 0, nil
	}

	num, err := strconv.ParseUint(numPart, 10, 64)
	if err != nil {
		return 0, err
	}

	return num * unitMultiplier, nil
}

func jvmArgsXMS(XMX uint64) uint64 {
	var bytes uint64
	var unit string
	for {
		options := []string{"大小", "单位", "确认"}
		selection := lib.Select(options, fmt.Sprintf("请选择一个选项, 当前XMS: %d%s", bytes, unit))
		switch selection {
		case 0:
			input := lib.Input("请输入XMS大小(数字, 不可小于0): ")
			num, err := strconv.ParseUint(input, 10, 64)
			if err != nil {
				fmt.Println("这不是一个有效的数字!")
			}
			bytes = num
		case 1:
			options := []string{
				"GB: 1000MB", "GiB: 1024MB",
				"MB: 1000KB", "MiB: 1024KB",
				"KB: 1024B", "KiB: 1024B",
				"B: 字节",
			}
			selection := lib.Select(options, "请选择一个单位")
			switch selection {
			case 0:
				unit = "GB"
			case 1:
				unit = "GiB"
			case 2:
				unit = "MB"
			case 3:
				unit = "MiB"
			case 4:
				unit = "KB"
			case 5:
				unit = "KiB"
			case 6:
				unit = "B"
			}
		case 2:
			XMS, err := ToBytes(fmt.Sprintf("%d%s", bytes, unit))
			if err != nil {
				panic(err)
			}
			if XMX != 0 && XMS > XMX {
				fmt.Println("XMS不能大于XMX")
				continue
			}
			if XMS <= JVM_min_ram {
				fmt.Println("XMS不能小于1MiB")
				continue
			}
			if RealMem < XMS {
				fmt.Println("XMS不能大于物理内存")
				continue
			}
			return XMS
		}
	}
}

func jvmArgsXMX(XMS uint64) uint64 {
	var bytes uint64 = 1
	unit := "GiB"
	for {
		options := []string{"大小", "单位", "确认"}
		selection := lib.Select(options, fmt.Sprintf("请选择一个选项, 当前XMX: %d%s", bytes, unit))
		switch selection {
		case 0:
			input := lib.Input("请输入XMX大小(数字, 不可小于0): ")
			num, err := strconv.ParseUint(input, 10, 64)
			if err != nil {
				fmt.Println("这不是一个有效的数字!")
			}
			bytes = num
		case 1:
			options := []string{
				"GB: 1000MB", "GiB: 1024MB",
				"MB: 1000KB", "MiB: 1024KB",
				"KB: 1024B", "KiB: 1024B",
				"B: 字节",
			}
			selection := lib.Select(options, "请选择一个单位")
			switch selection {
			case 0:
				unit = "GB"
			case 1:
				unit = "GiB"
			case 2:
				unit = "MB"
			case 3:
				unit = "MiB"
			case 4:
				unit = "KB"
			case 5:
				unit = "KiB"
			case 6:
				unit = "B"
			}
		case 2:
			XMX, err := ToBytes(fmt.Sprintf("%d%s", bytes, unit))
			if err != nil {
				panic(err)
			}
			if XMS != 0 && XMX < XMS {
				fmt.Println("XMX不能小于XMS")
				continue
			}
			if XMX <= JVM_min_ram {
				fmt.Println("XMX不能小于1MiB")
				continue
			}
			if RealMem < XMX {
				fmt.Println("XMX不能大于物理内存")
				continue
			}
			return XMX
		}
	}
}

func encoding() string {
	options := []string{"UTF-8", "GBK"}
	selection := lib.Select(options, "请选择一个编码")
	return options[selection]
}

func jvmArgs(jvmArgs []string) []string {
	for {
		options := jvmArgs
		options = append(options, "添加参数", "确认")
		selection := lib.Select(options, "请选择一个选项或要更改的Java虚拟机参数")
		switch selection {
		case len(options) - 2:
			input := lib.Input("请输入Java虚拟机参数: ")
			if input == "" {
				continue
			}
			jvmArgs = append(jvmArgs, input)
		case len(options) - 1:
			return jvmArgs
		default:
			input := lib.Input("请输入Java虚拟机参数, 为空即移除参数: ")
			if input == "" {
				for i, v := range jvmArgs {
					if v == options[selection] {
						jvmArgs = append(jvmArgs[:i], jvmArgs[i+1:]...)
					}
				}
			}
			jvmArgs[selection] = input
		}

	}
}

func java() lib.JavaInfo {
	for {
		javaLists, err := lib.LoadJavaLists()
		if err != nil {
			panic(err)
		}
		options := []string{}
		for _, v := range javaLists {
			options = append(options, fmt.Sprintf("%s(%s)", v.Version, v.Path))
		}
		options = append(options, "重新检测Java环境", "手动选择Java可执行程序")
		selection := lib.Select(options, "请选择一个Java环境或选项")

		switch selection {
		case len(options) - 2:
			javaLists = lib.DetectJava()
			err := lib.SaveJavaLists(javaLists)
			if err != nil {
				panic(err)
			}
			continue
		case len(options) - 1:
			input := lib.Input("请输入Java可执行程序路径: ")
			if _, err := os.Stat(input); err != nil && os.IsNotExist(err) {
				fmt.Println("路径不正确")
				continue
			}
			java_ver, err := lib.GetJavaVersion(input)
			if err != nil {
				fmt.Println("无法获取Java版本")
				continue
			}
			javaInfo := lib.JavaInfo{
				Version: java_ver,
				Path:    input,
			}
			javaLists = append(javaLists, javaInfo)
			err = lib.SaveJavaLists(javaLists)
			if err != nil {
				panic(err)
			}
			return javaInfo
		default:
			return javaLists[selection]
		}
	}
}

func serverArgs(serverArgs []string) []string {
	for {
		options := serverArgs
		options = append(options, "添加参数", "确认")
		selection := lib.Select(options, "请选择一个选项或要更改的Java虚拟机参数")
		switch selection {
		case len(options) - 2:
			input := lib.Input("请输入Java虚拟机参数: ")
			if input == "" {
				continue
			}
			serverArgs = append(serverArgs, input)
		case len(options) - 1:
			return serverArgs
		default:
			input := lib.Input("请输入Java虚拟机参数, 为空即移除参数: ")
			if input == "" {
				for i, v := range serverArgs {
					if v == options[selection] {
						serverArgs = append(serverArgs[:i], serverArgs[i+1:]...)
					}
				}
			}
			serverArgs[selection] = input
		}

	}
}
