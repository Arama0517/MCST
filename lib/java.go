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
		if strings.Contains(strings.ToLower(path), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

// searchFile 在指定目录下使用多线程搜索指定的文件名
func searchFile(path string, executePaths *[]string, executeName string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	for _, entry := range entries {
		if entry.IsDir() {
			wg.Add(1)
			go func(entry os.DirEntry) {
				defer wg.Done()
				if filepath.Base(entry.Name()) == executeName {
					fullPath := filepath.Join(path, entry.Name())
					*executePaths = append(*executePaths, fullPath)
				} else {
					if shouldExclude(entry.Name()) {
						return
					}
					searchFile(filepath.Join(path, entry.Name()), executePaths, executeName)
				}
			}(entry)
		} else {
			if filepath.Base(entry.Name()) == executeName {
				fullPath := filepath.Join(path, entry.Name())
				if !shouldExclude(fullPath) {
					*executePaths = append(*executePaths, fullPath)
				}
			}
		}
	}
	wg.Wait()
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

// detectJava 检测计算机上所有已安装的 Java
func DetectJava() []map[string]string {
	javaPaths := make([]string, 0)

	if runtime.GOOS == "windows" {
		for drive := 'A'; drive <= 'Z'; drive++ {
			root := string(drive) + ":\\"
			searchFile(root, &javaPaths, "java.exe")
		}
	} else {
		searchFile("/usr/lib", &javaPaths, "java")
	}

	javaWithVersion := make([]map[string]string, 0)
	for _, java := range javaPaths {
		version, err := getJavaVersion(java)
		if err == nil {
			javaWithVersion = append(javaWithVersion, map[string]string{"path": java, "version": version})
		}
	}

	return javaWithVersion
}