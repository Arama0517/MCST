/*
 * MCSCS can be used to easily create, launch, and configure a Minecraft server.
 * Copyright (C) 2024 Arama
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package lib

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var excludedKeywords = []string{"$", "{", "}", "_", "Windows", "AppData"}

func shouldExclude(path string) bool {
	for _, keyword := range excludedKeywords {
		if strings.Contains(path, keyword) {
			return true
		}
	}
	return false
}

// 控制并发数量
var semaphore = make(chan struct{}, 20)

// 搜索指定目录下的文件
func searchFile(path string, executeName string, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	semaphore <- struct{}{}        // 获取信号量
	defer func() { <-semaphore }() // 释放信号量

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if shouldExclude(entry.Name()) {
				continue
			}
			wg.Add(1)
			go searchFile(filepath.Join(path, entry.Name()), executeName, resultChan, wg)
		} else {
			if filepath.Base(entry.Name()) == executeName {
				fullPath := filepath.Join(path, entry.Name())
				if !shouldExclude(fullPath) {
					resultChan <- fullPath
				}
			}
		}
	}
}

// 获取 Java 版本信息
func GetJavaVersion(javaPath string) (string, error) {
	output, err := exec.Command(javaPath, "-version").CombinedOutput()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:[._](\d+))?(?:-(.+))?`)
	version := re.FindString(string(output))
	if version != "" {
		return version, nil
	}
	return "", fmt.Errorf("version not found")
}

type JavaInfo struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

// 检测已安装的 Java
func DetectJava() []JavaInfo {
	findJavas := []JavaInfo{}

	var wg sync.WaitGroup
	resultChan := make(chan string, 100)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	if runtime.GOOS == "windows" {
		for drive := 'A'; drive <= 'Z'; drive++ {
			root := string(drive) + ":\\"
			wg.Add(1)
			go searchFile(root, "java.exe", resultChan, &wg)
		}
	} else {
		wg.Add(1)
		go searchFile("/usr/lib", "java", resultChan, &wg)
	}

	for java := range resultChan {
		version, err := GetJavaVersion(java)
		if err == nil {
			findJavas = append(findJavas, JavaInfo{Path: java, Version: version})
		}
	}
	return findJavas
}
