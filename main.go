package main

import (
	"encoding/json"
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
	javas := lib.DetectJava()
	jsonJavas, err := json.Marshal(javas)
	if err != nil {
		Logger.Error(err)
		return
	}
	fmt.Println(string(jsonJavas))
	lib.KillAria2c()
}
