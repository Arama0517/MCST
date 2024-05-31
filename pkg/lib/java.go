/*
 * Minecraft Server Tool(MCST) is a command-line utility making Minecraft server creation quick and easy for beginners.
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
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"

	"github.com/charlievieth/fastwalk"
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
	return "", ErrJavaVersionNotFound
}

type JavaInfo struct {
	Path    string `json:"path"`
	Version string `json:"version"`
}

func searchFile(path string, name string) ([]string, error) {
	var results []string
	err := fastwalk.Walk(&fastwalk.DefaultConfig, path, func(path string, d fs.DirEntry, err error) error {
		if err != nil && !os.IsPermission(err) {
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

func detectJavaWindows(javaPaths *[]string) {
	for drive := 'A'; drive <= 'Z'; drive++ {
		root := string(drive) + ":\\"
		if _, err := os.Stat(root); err == nil {
			*javaPaths, err = searchFile(root, "java.exe")
			if err != nil {
				continue
			}
		}
	}
}

func DetectJava() ([]JavaInfo, error) {
	var findJavas []JavaInfo
	var javaPaths []string
	if runtime.GOOS == "windows" {
		detectJavaWindows(&javaPaths)
	} else {
		var err error
		javaPaths, err = searchFile("/usr/lib", "java")
		if err != nil {
			return nil, err
		}
	}
	var wg sync.WaitGroup
	for _, java := range javaPaths {
		wg.Add(1)
		go func(java string) {
			defer wg.Done()
			version, err := GetJavaVersion(java)
			if err == nil {
				findJavas = append(findJavas, JavaInfo{
					Path:    java,
					Version: version,
				})
			}
		}(java)
	}
	wg.Wait()
	return findJavas, nil
}
