package lib_test

import (
	"net/url"
	"os"
	"testing"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
)

var URL = url.URL{ // https://ash-speed.hetzner.com/100MB.bin
	Scheme: "https",
	Host:   "ash-speed.hetzner.com",
	Path:   "/100MB.bin",
}

func init() {
	if err := lib.InitAll(); err != nil {
		panic(err)
	}
}

func TestDownload(t *testing.T) {
	lib.EnableAria2c = false
	path, err := lib.NewDownloader(URL).Download()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}

func TestAria2Downlaod(t *testing.T) {
	lib.EnableAria2c = true
	path, err := lib.NewDownloader(URL).Download()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(path)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
}
