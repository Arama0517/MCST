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

package create

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Arama0517/MCST/internal/bytes"
	"github.com/Arama0517/MCST/internal/configs"
	MCSTErrors "github.com/Arama0517/MCST/internal/errors"
	"github.com/Arama0517/MCST/internal/locale"
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

func New() *cobra.Command {
	flags := createCmdFlags{}
	cmd := &cobra.Command{
		Use:               "create",
		Short:             locale.GetLocaleMessage("create.short"),
		Long:              locale.GetLocaleMessage("create.long"),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(_ *cobra.Command, _ []string) error {
			if !flags.eula {
				return MCSTErrors.ErrEulaRequired
			}
			var err error
			var config configs.Server
			config.Name = flags.name
			for name := range configs.Configs.Servers {
				if name == flags.name {
					return MCSTErrors.ErrServerExists
				}
			}
			config.Java.MinMemory, err = bytes.ToBytes(flags.xms)
			if err != nil {
				return err
			}
			config.Java.MaxMemory, err = bytes.ToBytes(flags.xmx)
			if err != nil {
				return err
			}
			memInfo, err := mem.VirtualMemory()
			if err != nil {
				return err
			}
			switch {
			case config.Java.MinMemory < bytes.MiB:
				return MCSTErrors.ErrXmsToLow
			case config.Java.MaxMemory < bytes.MiB:
				return MCSTErrors.ErrXmxTooLow
			case config.Java.MaxMemory < config.Java.MinMemory:
				return MCSTErrors.ErrXmxLessThanXms
			case config.Java.MinMemory > memInfo.Total:
				return MCSTErrors.ErrXmsExceedsPhysicalMemory
			case config.Java.MaxMemory > memInfo.Total:
				return MCSTErrors.ErrXmxExceedsPhysicalMemory
			}
			config.Java.Encoding = flags.encoding
			config.Java.Path = flags.java
			config.Java.Args = flags.jvmArgs
			config.ServerArgs = flags.serverArgs
			if flags.core < 0 || flags.core > len(configs.Configs.Cores) {
				return MCSTErrors.ErrCoreNotFound
			}
			configs.Configs.Servers[config.Name] = config

			// EULA 部分
			EULAFileData := fmt.Sprintf(`# Create By Minecraft Server Tool
# By changing the setting below to TRUE you are indicating your agreement to Minecraft EULA(<https://aka.ms/MinecraftEULA/>).
# %s
eula=true`, time.Now().Format("Mon Jan 02 15:04:05 MST 2006"))
			if err := os.MkdirAll(filepath.Join(configs.ServersDir, config.Name), 0o755); err != nil {
				return err
			}
			EULAFile, err := os.Open(filepath.Join(configs.ServersDir, config.Name, "eula.txt"))
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
			srcFile, err := os.Open(configs.Configs.Cores[flags.core].FilePath)
			if err != nil {
				return err
			}
			dstFile, err := os.Create(filepath.Join(configs.ServersDir, config.Name, "server.jar"))
			if err != nil {
				return err
			}
			if _, err = io.Copy(dstFile, srcFile); err != nil {
				return err
			}
			if err := os.Chmod(filepath.Join(configs.ServersDir, config.Name, "server.jar"), 0o755); err != nil {
				return err
			}
			if err := srcFile.Close(); err != nil {
				return err
			}
			if err := dstFile.Close(); err != nil {
				return err
			}
			if err := configs.Configs.Save(); err != nil {
				return err
			}
			log.Info("保存成功")
			return nil
		},
	}
	cmd.Flags().StringVarP(&flags.name, "name", "n", "", locale.GetLocaleMessage("create.flags.name"))
	cmd.Flags().StringVar(&flags.xms, "xms", "1G", locale.GetLocaleMessage("create.flags.xms"))
	cmd.Flags().StringVar(&flags.xmx, "xmx", "1G", locale.GetLocaleMessage("create.flags.xmx"))
	cmd.Flags().StringVarP(&flags.encoding, "encoding", "e", "UTF-8", locale.GetLocaleMessage("create.flags.encoding"))
	cmd.Flags().StringVarP(&flags.java, "java", "j", "", locale.GetLocaleMessage("create.flags.java"))
	cmd.Flags().StringSliceVar(&flags.jvmArgs, "jvm_args", []string{"-Dlog4j2.formatMsgNoLookups=true"}, locale.GetLocaleMessage("create.flags.jvm_args"))
	cmd.Flags().StringSliceVar(&flags.serverArgs, "server_args", []string{"--nogui"}, locale.GetLocaleMessage("create.flags.server_args"))
	cmd.Flags().IntVarP(&flags.core, "core", "c", 0, locale.GetLocaleMessage("create.flags.core"))
	if !configs.Configs.Settings.AutoAcceptEULA {
		cmd.Flags().BoolVar(&flags.eula, "eula", false, locale.GetLocaleMessage("create.flags.eula"))
		_ = cmd.MarkFlagRequired("eula")
	} else {
		flags.eula = true
	}
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("java")
	_ = cmd.MarkFlagRequired("core")
	return cmd
}
