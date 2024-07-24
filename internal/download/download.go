/*
 * Minecraft Server Tool(MCST) is a command-line utility making Minecraft server creation quick and easy for beginners.
 * Copyright (c) 2024-2024 Arama.
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

package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/requests"
	"github.com/apex/log"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	URL      string
	FileName string
	bar      *progressbar.ProgressBar
}

func NewDownloader(url string) *Downloader {
	return &Downloader{URL: url}
}

func (d *Downloader) Download() (string, error) {
	// 请求服务器
	req, err := requests.NewRequest(http.MethodGet, d.URL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	// 获取文件名
	d.FileName = resp.Header.Get("Content-Disposition")
	if d.FileName == "" {
		d.FileName = filepath.Base(req.URL.Path)
	}

	// 检测是否已经存在
	filePath := filepath.Join(configs.DownloadsDir, d.FileName)
	if _, err = os.Stat(filePath); err == nil {
		log.Info("检测到文件已存在, 已跳过下载")
		return filePath, nil
	}

	// 设置下载进度条
	d.bar = progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetDescription("[cyan]下载中...[reset]"),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[black][[reset]",
			BarEnd:        "[black]][reset]",
		}),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionThrottle(50*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			if _, err = fmt.Fprint(os.Stderr, "\n"); err != nil {
				return
			}
		}))

	// 下载
	if configs.Configs.Settings.Aria2.Enable {
		return d.aria2Download()
	}
	return d.defaultDownload(filePath, resp)
}

// defaultDownload 单线程下载
func (d *Downloader) defaultDownload(filePath string, resp *http.Response) (string, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	_, err = io.Copy(io.MultiWriter(file, d.bar), resp.Body)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
