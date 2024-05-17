package main

import (
	"flag"
	"fmt"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
	"github.com/Arama-Vanarana/MCSCS-Go/pages"
)

var Logger = lib.Logger

func main() {
	version := flag.Bool("version", false, "显示程序版本")
	flag.Parse()
	if *version {
		fmt.Println(lib.VERSION)
		return
	}
	options := []string{"创建服务器", "退出"}
	for {
		selected := lib.Select(options, "请选择一个选项 ")
		var err error
		switch selected {
		case 0:
			err = pages.CreatePage()
		default:
			return
		}
		if err != nil {
			Logger.WithError(err).Error("运行失败")
		}
	}
}
