package lib_test

import (
	"os"
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
	"github.com/sirupsen/logrus"
)

var testingFunc *testing.T

func TestLogs(t *testing.T) {
	testingFunc = t
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
	testingFunc.Log("logger.ExitFunc called, with exit code:", code)
}
