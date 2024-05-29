package lib

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/charlievieth/fastwalk"
	"github.com/rs/zerolog/log"
)

// https://github.com/MCSLTeam/MCSL2/blob/master/MCSL2Lib/ProgramControllers/javaDetector.py 第86行
var javaVersionRegex = regexp.MustCompile(`(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:[._](\d+))?(?:-(.+))?`)

// GetJavaVersion 获取 Java 的版本
func GetJavaVersion(javaPath string) (string, error) {
	output, err := exec.Command(javaPath, "-version").CombinedOutput()
	if err != nil {
		return "", err
	}
	matches := javaVersionRegex.FindStringSubmatch(string(output))
	if len(matches) > 0 {
		return matches[0], nil
	}
	return "", fmt.Errorf("failed to parse Java version from output")
}

type JavaInfo struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

func searchFile(path string, name string) ([]string, error) {
	var results []string
	err := fastwalk.Walk(&fastwalk.DefaultConfig, path, func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsPermission(err) {
			log.Error().Err(err).Msg("遍历Java目录失败")
			return nil
		}
		if ok, err := filepath.Match(name, d.Name()); !ok {
			return err
		}
		results = append(results, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

func DetectJava() ([]JavaInfo, error) {
	var findJavas []JavaInfo
	var javaPaths []string
	if runtime.GOOS == "windows" {
		for drive := 'A'; drive <= 'Z'; drive++ {
			root := string(drive) + ":\\"
			if _, err := os.Stat(root); err == nil {
				javaPaths, err = searchFile(root, "java.exe")
				if err != nil {
					continue
				}
			}
		}
	} else {
		var err error
		javaPaths, err = searchFile("/usr/lib", "java")
		if err != nil {
			return nil, err
		}
	}
	for _, java := range javaPaths {
		version, err := GetJavaVersion(java)
		if err == nil {
			findJavas = append(findJavas, JavaInfo{Path: java, Version: version})
		}
	}
	log.Info().Interface("javas", findJavas).Msg("寻找到的Java环境")
	return findJavas, nil
}
