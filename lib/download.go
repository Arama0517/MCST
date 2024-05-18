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
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// 下载文件, 返回文件路径
func Download(url string, fileName string) (string, error) {
	filePath := filepath.Join(DownloadsDir, fileName)
	// 不重复下载
	if _, err := os.Stat(filePath); err == nil {
		return filePath, nil
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "MCSCS-Golang/"+VERSION)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求%s失败, 状态码: %s", url, resp.Status)
	}

	// 进度条
	ansiStdout := ansi.NewAnsiStdout()
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("[cyan]下载中...[reset]"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			BarEnd:        "]",
			BarStart:      "[",
			Saucer:        "[green]━[reset]",
			SaucerPadding: "[red]━[reset]",
		}),
		progressbar.OptionSetWriter(ansiStdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(ansiStdout, "\n")
		}),
	)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
