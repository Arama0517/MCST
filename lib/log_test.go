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

package lib_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
	"github.com/sirupsen/logrus"
)

func TestLogs(t *testing.T) {
	logger := lib.Logger
	logger.ExitFunc = LoggerExitFunc
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
	logger.Trace("this is a trace message")
	logger.Debug("this is a debug message")
	logger.Info("this is an info message")
	logger.Warn("this is a warning message")
	logger.Error("this is an error message")
	logger.Fatal("this is a fatal message")
}

func LoggerExitFunc(code int) {
	fmt.Println("logger.ExitFunc called, with exit code:", code)
}
