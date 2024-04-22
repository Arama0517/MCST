package lib

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
)

func Select(options []string, message string) int {
	var result int
	prompt := &survey.Select{
		Options: options,
		Message: message,
	}
	err := survey.AskOne(prompt, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func SelectFile(filename string) (string, error) {
	// 初始目录为根目录
	var currentDir string
	if runtime.GOOS == "linux" {
		currentDir = "/"
	} else if runtime.GOOS == "windows" {
		currentDir = `C:\`
	}

main:
	for {

		// 构建选项列表
		options := make([]string, 0)
		// 获取当前目录下的所有文件和文件夹
		var files []fs.DirEntry
		if currentDir != "drives" {
			files_read, err := os.ReadDir(currentDir)
			if err != nil {
				log.Fatal(err)
			}
			files = files_read
		}
		// 如果不是根目录，添加 ".." 返回上一级目录的选项
		if currentDir != "/" && currentDir != "drives" {
			options = append(options, "..")
		}
		if currentDir == "drives" {
			options = append(options, getDrivePaths()...)
		} else {
			for _, file := range files {
				options = append(options, file.Name())
			}
		}

		// 询问用户选择文件或文件夹
		selectedIndex := Select(options, "请选择一个文件或文件夹")

		selectedName := options[selectedIndex]
		selectedPath := filepath.Join(currentDir, selectedName)

		// 如果用户选择的是文件夹，则进入该文件夹，否则检查文件名称是否与指定名称一致
		if fileInfo, err := os.Stat(selectedPath); err == nil && fileInfo.IsDir() {
			if selectedName != ".." {
				currentDir = selectedPath
			} else {
				for _, drivePath := range getDrivePaths() {
					if currentDir == drivePath {
						currentDir = "drives"
						continue main
					}
				}
				currentDir = filepath.Dir(currentDir)
			}
		} else if currentDir == "drives" {
			currentDir = selectedPath
		} else {
			// 检查文件名称是否与指定名称一致
			if selectedName != filename {
				fmt.Println("请选择名为 '", filename, "' 的文件")
				continue
			}
			return selectedPath, nil
		}
	}
}

func getDrivePaths() []string {
	var drivePaths []string
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drivePath := string(drive) + `:\`
		if _, err := os.Stat(drivePath); err == nil {
			drivePaths = append(drivePaths, drivePath)
		}
	}
	return drivePaths
}
