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
		lib.KillAria2c()
		return
	}
	fmt.Println(lib.Select([]string{"Hello", "World", "Go"}, "aaa"))
	lib.KillAria2c()
}
