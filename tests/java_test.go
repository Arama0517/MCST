package main

import (
	"fmt"

	"github.com/Arama-Vanarana/MCSCS-Go/lib"
)

func TestDetectJava() {
	javaList := lib.DetectJava()
	for _, java := range javaList {
		fmt.Printf("Path: %s, Version: %s\n", java["path"], java["version"])
	}
}
