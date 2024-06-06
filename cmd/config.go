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
	"errors"
	"os"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/caarlos0/log"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
)

type configCmdFlags struct {
	name       string
	xms        string
	xmx        string
	encoding   string
	java       string
	jvmArgs    []string
	serverArgs []string
}

func newConfigCmd() *cobra.Command {
	flags := &configCmdFlags{}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置服务器",
		Long:  "配置服务器的各项参数",
		RunE: func(cmd *cobra.Command, _ []string) error {
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			config, exists := configs.Servers[flags.name]
			if !exists {
				return errors.New("服务器不存在")
			}
			cmdFlags := cmd.Flags()
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				return err
			}
			if cmdFlags.Changed("xms") {
				ram, err := toBytes(flags.xms)
				if err != nil {
					return err
				}
				switch {
				case ram > config.Java.Xmx:
					return errors.New("参数错误 Xms 不可大于 Xmx")
				case ram < 1048576:
					return errors.New("参数错误 Xms 不可小于 1MiB")
				case ram > memInfo.Total:
					return errors.New("参数错误 Xms 不可大于物理内存")
				}
				config.Java.Xms = ram
			}
			if cmdFlags.Changed("xmx") {
				ram, err := toBytes(flags.xmx)
				if err != nil {
					return err
				}
				switch {
				case ram < config.Java.Xms:
					return errors.New("参数错误 Xmx 不可小于 Xms")
				case ram < 1048576:
					return errors.New("参数错误 Xmx 不可小于 1MiB")
				case ram > memInfo.Total:
					return errors.New("参数错误 Xmx 不可大于物理内存")
				}
				config.Java.Xmx = ram
			}
			if cmdFlags.Changed("encoding") {
				config.Java.Encoding = flags.encoding
			}
			if cmdFlags.Changed("java") {
				if _, err := os.Stat(flags.java); err != nil {
					return err
				} else {
					config.Java.Path = flags.java
				}
			}
			if cmdFlags.Changed("jvm_args") {
				config.Java.Args = flags.jvmArgs
			}
			if cmdFlags.Changed("server_args") {
				config.Java.Args = flags.serverArgs
			}
			log.Info("设置完成")
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.name, "name", "n", "", "服务器名称")
	cmd.Flags().StringVar(&flags.xms, "xms", "", "Java虚拟机初始堆内存")
	cmd.Flags().StringVar(&flags.xmx, "xmx", "", "Java虚拟机最大堆内存")
	cmd.Flags().StringVarP(&flags.encoding, "encoding", "e", "", "输出编码")
	cmd.Flags().StringVarP(&flags.java, "java", "j", "", "使用的Java")
	cmd.Flags().StringSliceVar(&flags.jvmArgs, "jvm_args", []string{}, "Java虚拟机其他参数")
	cmd.Flags().StringSliceVar(&flags.serverArgs, "server_args", []string{}, "Minecraft服务器参数")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
