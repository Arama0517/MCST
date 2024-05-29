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

package lib

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	LogFile, err := os.Create(filepath.Join(logsDir, time.Now().Format("2006010215")+".log"))
	if err != nil {
		panic(err)
	}
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	log.Logger = zerolog.New(LogFile).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	log.Info().Msg("Hello world!")
	log.Info().Msg("本程序遵循GPLv3协议开源")
	log.Info().Msg("作者: Arama 3584075812@qq.com")
	log.Info().Msgf("MCSCS-Go 版本: %s", Version)
	log.Info().Str("GOVERSION", runtime.Version()).Str("GOOS", runtime.GOOS).Str("GOARCH", runtime.GOARCH).Send()
	configs, err := LoadConfigs()
	if err != nil {
		panic(err)
	}
	LogLevel, err := zerolog.ParseLevel(configs.LogLevel)
	if err != nil {
		zerolog.SetGlobalLevel(LogLevel)
	}
}
