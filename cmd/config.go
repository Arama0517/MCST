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
	"os"
	"path/filepath"

	"github.com/Arama-Vanarana/MCServerTool/internal/lib"
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
	delete     bool // 如果为true就删除服务器
}

func newConfigCmd() *cobra.Command {
	flags := &configCmdFlags{}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "配置服务器",
		Long: `配置服务器的各项参数
如果使用 '--delete' 参数则删除服务器(慎重, 无法复原)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			configs, err := lib.LoadConfigs()
			if err != nil {
				return err
			}
			config, exists := configs.Servers[flags.name]
			if !exists {
				return ErrServerNotFound
			}
			cmdFlags := cmd.Flags()
			if flags.delete {
				delete(configs.Servers, flags.name)
				if err := os.RemoveAll(filepath.Join(lib.ServersDir, flags.name)); err != nil {
					return err
				}
				log.Info("删除服务器成功")
				return nil
			}
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
					return ErrXmxLessThanXms
				case ram < 1048576:
					return ErrXmsToLow
				case ram > memInfo.Total:
					return ErrXmsExceedsPhysicalMemory
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
					return ErrXmxLessThanXms
				case ram < 1048576:
					return ErrXmxTooLow
				case ram > memInfo.Total:
					return ErrXmxExceedsPhysicalMemory
				}
				config.Java.Xmx = ram
			}
			if cmdFlags.Changed("encoding") {
				config.Java.Encoding = flags.encoding
			}
			if cmdFlags.Changed("java") {
				if _, err := os.Stat(flags.java); err != nil {
					return err
				}
				config.Java.Path = flags.java
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
	cmd.Flags().BoolVar(&flags.delete, "delete", false, "删除服务器")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}