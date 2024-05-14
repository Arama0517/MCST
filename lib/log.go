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
	logFile, err := os.OpenFile(filepath.Join(GetLogsDir(), "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
