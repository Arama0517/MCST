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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Arama0517/MCST/pkg/lib"
	"github.com/apex/log"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

type createCmdFlags struct {
	name       string
	xms        string
	xmx        string
	encoding   string
	java       string
	jvmArgs    []string
	serverArgs []string
	core       int
	eula       bool
}

func newCreateCmd() *cobra.Command {
	flags := createCmdFlags{}
	cmd := &cobra.Command{
		Use:   "create",
		Short: "创建服务器",
		Long: `如果你还未下载任何核心, 请使用 'MCST download' 下载核心
必须指定--name, --java, --core`,
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(_ *cobra.Command, _ []string) error {
			if !flags.eula {
				return ErrEulaRequired
			}
			var err error
			var config lib.Server
			config.Name = flags.name
			for name := range lib.Configs.Servers {
				if name == flags.name {
					return ErrServerExists
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
			case config.Java.Xms < MiB:
				return ErrXmsToLow
			case config.Java.Xmx < MiB:
				return ErrXmxTooLow
			case config.Java.Xmx < config.Java.Xms:
				return ErrXmxLessThanXms
			case config.Java.Xms > memInfo.Total:
				return ErrXmsExceedsPhysicalMemory
			case config.Java.Xmx > memInfo.Total:
				return ErrXmxExceedsPhysicalMemory
			}
			config.Java.Encoding = flags.encoding
			config.Java.Path = flags.java
			config.Java.Args = flags.jvmArgs
			config.ServerArgs = flags.serverArgs
			if flags.core < 0 || flags.core > len(lib.Configs.Cores) {
				return ErrCoreNotFound
			}
			lib.Configs.Servers[config.Name] = config

			// EULA 部分
			EULAFileData := fmt.Sprintf(`# Create By Minecraft Server Tool
# By changing the setting below to TRUE you are indicating your agreement to Minecraft EULA(<https://aka.ms/MinecraftEULA/>).
# %s
eula=true`, time.Now().Format("Mon Jan 02 15:04:05 MST 2006"))
			if err := os.MkdirAll(filepath.Join(lib.ServersDir, config.Name), 0o755); err != nil {
				return err
			}
			EULAFile, err := os.Open(filepath.Join(lib.ServersDir, config.Name, "eula.txt"))
			if err != nil {
				return err
			}
			if _, err := EULAFile.WriteString(EULAFileData); err != nil {
				return err
			}
			if err := EULAFile.Close(); err != nil {
				return err
			}

			// 保存
			srcFile, err := os.Open(lib.Configs.Cores[flags.core].FilePath)
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
			if err := os.Chmod(filepath.Join(lib.ServersDir, config.Name, "server.jar"), 0o755); err != nil {
				return err
			}
			if err := srcFile.Close(); err != nil {
				return err
			}
			if err := dstFile.Close(); err != nil {
				return err
			}
			if err := lib.Configs.Save(); err != nil {
				return err
			}
			log.Info("保存成功")
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.name, "name", "n", "", "服务器名称")
	cmd.Flags().StringVar(&flags.xms, "xms", "1G", "Java虚拟机初始堆内存")
	cmd.Flags().StringVar(&flags.xmx, "xmx", "1G", "Java虚拟机最大堆内存")
	cmd.Flags().StringVarP(&flags.encoding, "encoding", "e", "UTF-8", "输出编码")
	cmd.Flags().StringVarP(&flags.java, "java", "j", "", "使用的Java")
	cmd.Flags().StringSliceVar(&flags.jvmArgs, "jvm_args", []string{"-Dlog4j2.formatMsgNoLookups=true"}, "Java虚拟机其他参数")
	cmd.Flags().StringSliceVar(&flags.serverArgs, "server_args", []string{"--nogui"}, "Minecraft服务器参数")
	cmd.Flags().IntVarP(&flags.core, "core", "c", 0, "使用的核心ID")
	cmd.Flags().BoolVar(&flags.eula, "eula", false, "是否同意EULA协议(https://aka.ms/MinecraftEULA/)")
	if !lib.Configs.AutoAcceptEULA {
		_ = cmd.MarkFlagRequired("eula")
	} else {
		flags.eula = true
	}
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("java")
	_ = cmd.MarkFlagRequired("core")
	return cmd
}

var (
	TiB uint64 = 1099511627776 // Tebibyte: 1024 * 1024 * 1024 * 1024
	TB  uint64 = 1000000000000 // Terabyte: 1000 * 1000 * 1000 * 1000
	GiB uint64 = 1073741824    // Gibibyte: 1024 * 1024 * 1024
	GB  uint64 = 1000000000    // Gigabyte: 1000 * 1000 * 1000
	MiB uint64 = 1048576       // Mebibyte: 1024 * 1024
	MB  uint64 = 1000000       // Megabyte: 1000 * 1000
	KiB uint64 = 1024          // Kibibyte
	KB  uint64 = 1000          // Kilobyte
)

func toBytes(byteStr string) (uint64, error) {
	var num, unit string

	// 分离数字部分和单位部分
	for _, char := range byteStr {
		switch {
		case char >= '0' && char <= '9':
			num += string(char)
		case (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z'):
			unit += strings.ToUpper(string(char))
		default:
			return 0, nil
		}
	}

	parsedNum, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}
	switch unit {
	case "T", "TIB":
		return parsedNum * TiB, nil
	case "TB":
		return parsedNum * TB, nil
	case "G", "GIB":
		return parsedNum * GiB, nil
	case "GB":
		return parsedNum * GB, nil
	case "M", "MIB":
		return parsedNum * MiB, nil
	case "MB":
		return parsedNum * MB, nil
	case "K", "KIB":
		return parsedNum * KiB, nil
	case "KB":
		return parsedNum * KB, nil
	case "B", "BYTES", "":
		return parsedNum, nil
	default:
		return 0, ErrInvalidUnit
	}
}
