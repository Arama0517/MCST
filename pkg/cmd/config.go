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

	"github.com/Arama0517/MCST/internal/bytes"
	"github.com/Arama0517/MCST/internal/configs"
	MCSTErrors "github.com/Arama0517/MCST/internal/errors"
	"github.com/Arama0517/MCST/internal/locale"
	"github.com/apex/log"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		Use:               "config",
		Short:             locale.GetLocaleMessage("config.short"),
		Long:              locale.GetLocaleMessage("config.long"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, _ []string) error {
			isCheckConfig := true
			cmdFlags := cmd.Flags()
			cmdFlags.VisitAll(func(flag *pflag.Flag) {
				if flag.Changed {
					isCheckConfig = false
				}
			})
			config, exists := configs.Configs.Servers[flags.name]
			if !exists {
				return MCSTErrors.ErrServerNotFound
			}
			if isCheckConfig {
				result := make(map[string]any)
				structToMap(config, "", result)
				for k, v := range result {
					log.WithField("value", v).Info(k)
				}
				return nil
			}
			if flags.delete {
				delete(configs.Configs.Servers, flags.name)
				if err := os.RemoveAll(filepath.Join(configs.ServersDir, flags.name)); err != nil {
					return err
				}
				return nil
			}
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				return err
			}
			if cmdFlags.Changed("xms") {
				ram, err := bytes.ToBytes(flags.xms)
				if err != nil {
					return err
				}
				switch {
				case ram > config.Java.MaxMemory:
					return MCSTErrors.ErrXmxLessThanXms
				case ram < bytes.MiB:
					return MCSTErrors.ErrXmsToLow
				case ram > memInfo.Total:
					return MCSTErrors.ErrXmsExceedsPhysicalMemory
				}
				config.Java.MinMemory = ram
			}
			if cmdFlags.Changed("xmx") {
				ram, err := bytes.ToBytes(flags.xmx)
				if err != nil {
					return err
				}
				switch {
				case ram < config.Java.MinMemory:
					return MCSTErrors.ErrXmxLessThanXms
				case ram < bytes.MiB:
					return MCSTErrors.ErrXmxTooLow
				case ram > memInfo.Total:
					return MCSTErrors.ErrXmxExceedsPhysicalMemory
				}
				config.Java.MaxMemory = ram
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
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.name, "name", "n", "", locale.GetLocaleMessage("config.flags.name"))
	cmd.Flags().StringVar(&flags.xms, "xms", "", locale.GetLocaleMessage("create.flags.xms"))
	cmd.Flags().StringVar(&flags.xmx, "xmx", "", locale.GetLocaleMessage("create.flags.xmx"))
	cmd.Flags().StringVarP(&flags.encoding, "encoding", "e", "", locale.GetLocaleMessage("create.flags.encoding"))
	cmd.Flags().StringVarP(&flags.java, "java", "j", "", locale.GetLocaleMessage("create.flags.java"))
	cmd.Flags().StringSliceVar(&flags.jvmArgs, "jvm_args", []string{}, locale.GetLocaleMessage("create.flags.jvm_args"))
	cmd.Flags().StringSliceVar(&flags.serverArgs, "server_args", []string{}, locale.GetLocaleMessage("create.flags.server_args"))
	cmd.Flags().BoolVar(&flags.delete, "delete", false, locale.GetLocaleMessage("config.flags.delete"))
	_ = cmd.MarkFlagRequired("name")
	return cmd
}
