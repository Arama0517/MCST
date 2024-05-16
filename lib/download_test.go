package lib_test

import (
	"testing"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

func TestDownload(t *testing.T) {
	lib.Download("https://dl.google.com/go/go1.14.2.src.tar.gz", "go1.14.2.src.tar.gz")
}
