package main

import (
	"flag"
	"fmt"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

var Logger = lib.Logger

func main() {
	version := flag.Bool("version", false, "显示程序版本")
	flag.Parse()
	if *version {
		fmt.Println(lib.VERSION)
		return
	}
	lib.Download("https://dl.google.com/go/go1.14.2.src.tar.gz", "go1.14.2.src.tar.gz")
}
