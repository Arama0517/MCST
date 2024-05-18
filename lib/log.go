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
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type MCSCSLogger struct {
	*logrus.Logger
}

var Logger = MCSCSLogger{logrus.New()}

func initLogs() {
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&logrus.TextFormatter{})
	logFile, err := os.OpenFile(filepath.Join(LogsDir, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Logger.SetOutput(logFile)
	Logger.Info("Hello world!")
	Logger.Info("本程序遵循GPLv3协议开源")
	Logger.Info("作者: Arama 3584075812@qq.com")
}

func (Logger *MCSCSLogger) Shell(cmd *exec.Cmd) {
	Logger.WithFields(logrus.Fields{
		"cmd": cmd,
	}).Info("运行命令")
}
