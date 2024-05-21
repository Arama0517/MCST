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
	"io/fs"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
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

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}

func SelectFile(filename string) (string, error) {
	// 初始目录为根目录
	currentDir, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
main:
	for {
		ClearScreen()
		options := make([]string, 0)
		displayOptions := make([]string, 0)
		var files []fs.DirEntry
		if currentDir != "drives" {
			files_read, err := os.ReadDir(currentDir)
			if err != nil {
				return "", err
			}
			files = files_read
		}
		if currentDir == "drives" {
			options = append(options, getDrivePaths()...)
		} else {
			if currentDir != "/" {
				options = append(options, "..")
				displayOptions = append(displayOptions, "返回上一级目录")
			}
			for _, file := range files {
				name := file.Name()
				displayName := name
				if file.IsDir() {
					displayName = "\033[32m" + name + "/\033[0m"
				}
				options = append(options, name)
				displayOptions = append(displayOptions, displayName)
			}
		}
		var displayCurrentDir string
		if currentDir == "drives" {
			displayCurrentDir = "分区选择"
		} else {
			displayCurrentDir = currentDir
		}
		selectedIndex := Select(displayOptions, "请选择文件: "+filename+", 当前路径: "+displayCurrentDir)
		selectedName := options[selectedIndex]
		selectedPath := filepath.Join(currentDir, selectedName)
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
			currentDir = selectedName
		} else {
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

func Input(message string) string {
	var result string
	prompt := &survey.Input{
		Message: message,
	}
	survey.AskOne(prompt, &result)
	return result
}

func Confirm(message string) bool {
	var result bool
	prompt := &survey.Confirm{
		Message: message,
	}
	survey.AskOne(prompt, &result)
	return result
}