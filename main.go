package main

import (
	"fmt"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
	log "github.com/sirupsen/logrus"
)

const (
	Version = "0.0.1"
)

func main() {
	log.SetReportCaller(true)
	result, err := lib.SelectFile("java")
	if err != nil {
		panic(err)
	}
	fmt.Println("Java: ", result)
}
