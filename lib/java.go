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

// searchFile 在指定目录下使用多线程搜索指定的文件名
func searchFile(path string, executeName string, resultChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
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

// getJavaVersion 获取 Java 的版本
func getJavaVersion(javaPath string) (string, error) {
	output, err := exec.Command(javaPath, "-version").CombinedOutput()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:[._](\d+))?(?:-(.+))?`)
	matches := re.FindStringSubmatch(string(output))
	if len(matches) > 0 {
		return matches[0], nil
	}
	return "", fmt.Errorf("regex")
}

type JavaInfo struct {
	Path    string
	Version string
}

// DetectJava 检测计算机上所有已安装的 Java
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
		version, err := getJavaVersion(java)
		if err == nil {
			findJavas = append(findJavas, JavaInfo{Path: java, Version: version})
		}
	}
	return findJavas
}
