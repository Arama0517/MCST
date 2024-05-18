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
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var currentTime = time.Now()

var DataDir string
var ConfigsDir string
var ServersDir string
var DownloadsDir string
var LogsDir string

func initDataDirs() {
	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DataDir = filepath.Join(homePath, ".config", "MCSCS")
	year, month, day := currentTime.Date()
	ConfigsDir = filepath.Join(DataDir, "configs")
	ServersDir = filepath.Join(DataDir, "servers")
	DownloadsDir = filepath.Join(DataDir, "downloads")
	LogsDir = filepath.Join(DataDir, "logs", fmt.Sprintf("%d%02d%02d%02d", year, month, day, currentTime.Hour()))
	createDirIfNotExist(DataDir)
	createDirIfNotExist(ConfigsDir)
	createDirIfNotExist(ServersDir)
	createDirIfNotExist(DownloadsDir)
	createDirIfNotExist(LogsDir)
}

func createDirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			panic(err)
		}
	}
}
