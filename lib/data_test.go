package lib_test

import (
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

func TestGetDatasDirs(t *testing.T) {
	t.Log("Data root directory: ", lib.DataDir)
	t.Log("Configs directory: ", lib.ConfigsDir)
	t.Log("Servers directory: ", lib.ServersDir)
	t.Log("Downloads directory: ", lib.DownloadsDir)
	t.Log("Logs directory: ", lib.LogsDir)
}
