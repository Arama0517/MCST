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
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/urfave/cli/v2"
)

var Create = cli.Command{
	Name:  "create",
	Usage: "创建服务器",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "服务器名称",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "Xms",
			Aliases: []string{"m"},
			Usage:   "Xms, Java虚拟机初始堆内存(可用单位: B, KB, KiB, MB, MiB, GB, GiB, TB, TiB), 默认为1G",
			Value:   "1G",
		},
		&cli.StringFlag{
			Name:    "Xmx",
			Aliases: []string{"x"},
			Usage:   "Xmx, Java虚拟机最大堆内存(可用单位: B, KB, KiB, MB, MiB, GB, GiB, TB, TiB); 默认为1G",
			Value:   "1G",
		},
		&cli.StringFlag{
			Name:    "encoding",
			Aliases: []string{"e"},
			Usage:   "使用的编码格式, 常见的编码格式: UTF-8, GBK, ASCII",
			Value:   "UTF-8",
		},
		&cli.PathFlag{
			Name:     "java",
			Aliases:  []string{"j"},
			Usage:    "Java路径",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:    "jvm_args",
			Aliases: []string{"a"},
			Usage:   "Java虚拟机的参数",
			Value:   cli.NewStringSlice("-Dlog4j2.formatMsgNoLookups=true"),
		},
		&cli.IntFlag{
			Name:     "core",
			Aliases:  []string{"c"},
			Usage:    "服务器核心, 从download命令下载, 使用downloads --list查看已下载的核心",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:    "server_args",
			Aliases: []string{"s"},
			Usage:   "Minecraft服务器特有参数",
			Value:   cli.NewStringSlice("--nogui"),
		},
	},
	Action: create,
}

func create(ctx *cli.Context) error {
	// 配置
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
	var config lib.Server
	config.Name = ctx.String("name")
	for name := range configs.Servers {
		if name == config.Name {
			return errors.New("服务器已存在")
		}
	}
	config.Java.Xms, err = toBytes(ctx.String("Xms"))
	if err != nil {
		return err
	}
	config.Java.Xmx, err = toBytes(ctx.String("Xmx"))
	if err != nil {
		return err
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	switch {
	case config.Java.Xms > config.Java.Xmx:
		return errors.New("Xms不能大于Xmx")
	case config.Java.Xms < 1024*1024, config.Java.Xmx < 1024*1024:
		return errors.New("Xms和Xmx必须大于1MiB")
	case config.Java.Xms > memInfo.Total, config.Java.Xmx > memInfo.Total:
		return errors.New("Xms和Xmx不能大于系统内存")
	}
	config.Java.Encoding = "UTF-8"
	if ctx.Bool("gbk") {
		config.Java.Encoding = "GBK"
	}
	config.Java.Path = ctx.Path("java")
	config.Java.Args = ctx.StringSlice("jvm_args")
	config.ServerArgs = ctx.StringSlice("server_args")
	coreIndex := ctx.Int("core")
	if coreIndex < 0 || coreIndex >= len(configs.Cores) {
		return errors.New("核心不存在")
	}
	coreInfo := configs.Cores[coreIndex]
	configs.Servers[config.Name] = config
	// EULA
	serverPath := filepath.Join(lib.ServersDir, config.Name)
	if err := os.MkdirAll(serverPath, 0755); err != nil {
		return err
	}
	choice, err := confirm("你是否同意EULA协议<https://aka.ms/MinecraftEULA/>?")
	if err != nil {
		return err
	}
	if choice {
		file, err := os.Create(filepath.Join(serverPath, "eula.txt"))
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)
		data := fmt.Sprintf(`# Create By Minecraft Server Tool
# By changing the setting below to TRUE you are indicating your agreement to Minecraft EULA(<https://aka.ms/MinecraftEULA/>).
# %s
eula=true`, time.Now().Format("Mon Jan 02 15:04:05 MCST 2006"))
		_, err = file.WriteString(data)
		if err != nil {
			return err
		}
	} else {
		return errors.New("你必须同意EULA协议<https://aka.ms/MinecraftEULA/>才能创建服务器")
	}
	// 保存
	if err := configs.Save(); err != nil {
		return err
	}
	// 复制
	srcFile, err := os.Open(coreInfo.FilePath)
	if err != nil {
		return err
	}
	dstFile, err := os.Create(filepath.Join(serverPath, "server.jar"))
	if err != nil {
		return err
	}
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if err := os.Chmod(filepath.Join(serverPath, "server.jar"), 0755); err != nil {
		return err
	}
	if err := srcFile.Close(); err != nil {
		return err
	}
	if err := dstFile.Close(); err != nil {
		return err
	}
	return nil
}

var ErrorNotAUnit = errors.New("这不是一个有效的单位")

func toBytes(byteStr string) (uint64, error) {
	var num, unit string

	// 分离数字部分和单位部分
	for _, char := range byteStr {
		if char >= '0' && char <= '9' {
			num += string(char)
		} else if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
			unit += strings.ToUpper(string(char))
		} else {
			return 0, nil
		}
	}

	var multiplier uint64
	switch unit {
	case "T", "TB", "TIB":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "G", "GB", "GIB":
		multiplier = 1024 * 1024 * 1024
	case "M", "MB", "MIB":
		multiplier = 1024 * 1024
	case "K", "KB", "KIB":
		multiplier = 1024
	case "B", "BYTES", "":
		multiplier = 1
	default:
		return 0, ErrorNotAUnit
	}

	parsedNum, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}

	return parsedNum * multiplier, nil
}

func confirm(description string) (bool, error) {
	fmt.Printf("%s (y/n): ", description)
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	choice = strings.TrimSpace(strings.ToLower(choice))
	for {
		switch choice {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			continue
		}
	}
}
