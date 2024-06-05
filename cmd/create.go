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

package cmd

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
	"github.com/spf13/cobra"
)

type createCmdFlage struct {
	name       string
	xms        string
	xmx        string
	encoding   string
	java       string
	jvmArgs    []string
	serverArgs []string
	core       int
}

func newCreateCmd() (*cobra.Command, error) {
	flags := createCmdFlage{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建服务器",
		Long: `如果你还未下载任何核心, 请使用 'MCST download' 下载核心
必须指定--name, --java, --core`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return create(flags, cmd)
		},
	}
	cmd.Flags().StringVarP(&flags.name, "name", "n", "", "服务器名称")
	cmd.Flags().StringVar(&flags.xms, "xms", "1G", "Java虚拟机初始堆内存")
	cmd.Flags().StringVar(&flags.xmx, "xmx", "1G", "Java虚拟机最大堆内存")
	cmd.Flags().StringVarP(&flags.encoding, "encoding", "e", "1G", "输出编码")
	cmd.Flags().StringVarP(&flags.java, "java", "j", "", "使用的Java")
	cmd.Flags().StringSliceVar(&flags.jvmArgs, "jvm_args", []string{"-Dlog4j2.formatMsgNoLookups=true"}, "Java虚拟机其他参数")
	cmd.Flags().StringSliceVar(&flags.serverArgs, "server_args", []string{"--nogui"}, "Minecraft服务器参数")
	cmd.Flags().IntVarP(&flags.core, "core", "c", 0, "使用的核心")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("java"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("core"); err != nil {
		return nil, err
	}
	return cmd, nil
}

func create(flags createCmdFlage, cmd *cobra.Command) error {
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
	var config lib.Server
	config.Name = flags.name
	for name := range configs.Servers {
		if name == flags.name {
			return errors.New("参数错误 服务器已存在, 请更换一个服务器名称")
		}
	}
	config.Java.Xms, err = toBytes(flags.xms)
	if err != nil {
		return err
	}
	config.Java.Xmx, err = toBytes(flags.xmx)
	if err != nil {
		return err
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	switch {
	case config.Java.Xmx < config.Java.Xms:
		return errors.New("参数错误 Xmx 不可小于 Xms")
	case config.Java.Xmx < 1048576, config.Java.Xms < 1048576:
		return errors.New("参数错误 Xmx 或 Xms 不可小于 1MiB")
	case config.Java.Xmx > memInfo.Total, config.Java.Xms > memInfo.Total:
		return errors.New("参数错误 Xmx 或 Xms 不可大于物理内存")
	}
	config.Java.Encoding = flags.encoding
	config.Java.Path = flags.java
	config.Java.Args = flags.jvmArgs
	config.ServerArgs = flags.serverArgs
	if flags.core < 0 || flags.core > len(configs.Cores) {
		return errors.New("参数错误 核心不存在")
	}
	configs.Servers[config.Name] = config

	// EULA 部分
	EULAFileData := fmt.Sprintf(`# Create By Minecraft Server Tool
# By changing the setting below to TRUE you are indicating your agreement to Minecraft EULA(<https://aka.ms/MinecraftEULA/>).
# %s
eula=true`, time.Now().Format("Mon Jan 02 15:04:05 MCST 2006"))
	if !configs.AutoAcceptEULA {
		if choice, err := confirm("你是否同意EULA协议<https://aka.ms/MinecraftEULA/>?"); err != nil {
			return err
		} else if !choice {
			return errors.New("你必须同意EULA协议<https://aka.ms/MinecraftEULA/>才能创建服务器")
		}
	}
	if err := configs.Save(); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(lib.ServersDir, config.Name), 0755); err != nil {
		return err
	}
	if EULAFile, err := os.Open(filepath.Join(lib.ServersDir, config.Name, "eula.txt")); err != nil {
		return err
	} else {
		if _, err := EULAFile.WriteString(EULAFileData); err != nil {
			return err
		}
		if err := EULAFile.Close(); err != nil {
			return err
		}
	}

	// 保存
	if err := configs.Save(); err != nil {
		return err
	}
	srcFile, err := os.Open(configs.Cores[flags.core].FilePath)
	if err != nil {
		return err
	}
	dstFile, err := os.Create(filepath.Join(lib.ServersDir, config.Name, "server.jar"))
	if err != nil {
		return err
	}
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if err := os.Chmod(filepath.Join(lib.ServersDir, config.Name, "server.jar"), 0755); err != nil {
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
		return 0, errors.New("这不是一个有效的单位")
	}

	parsedNum, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}

	return parsedNum * multiplier, nil
}

func confirm(description string) (bool, error) {
	if _, err := fmt.Printf("%s (y/n): ", description); err != nil {
		return false, err
	}
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
