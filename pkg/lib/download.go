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
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	URL            url.URL                  // 下载的 URL
	Multi          int                      // 是否多线程下载: 0 自动选择(默认), 1 单线程下载, 2 多线程下载
	fileName       string                   // 文件名
	bar            *progressbar.ProgressBar // 进度条
	maxConnections int64                    // 最大连接数
	contentLength  int64                    // 内容长度
}

func (d *Downloader) Download() (string, error) {
	resp, err := Request(d.URL, http.MethodHead, nil)
	if err != nil {
		return "", err
	}

	d.getFileName(resp.Header, d.URL)
	path := filepath.Join(DownloadsDir, d.fileName)
	fmt.Println(path)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	ansiStdout := ansi.NewAnsiStdout()
	d.bar = progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
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
		}))
	configs, err := LoadConfigs()
	if err != nil {
		return path, err
	}
	d.maxConnections = int64(configs.MaxConnections)
	d.contentLength = resp.ContentLength
	switch d.Multi {
	case 1:
		return path, d.singleDownload()
	case 2:
		if resp.Header.Get("Accept-Ranges") != "bytes" {
			return path, d.singleDownload()
		}
		return path, d.multiDownload()
	default:
		if resp.Header.Get("Accept-Ranges") == "bytes" && resp.ContentLength > 1024*1024 {
			return path, d.multiDownload()
		}
		return path, d.singleDownload()
	}
}

func (d *Downloader) singleDownload() error {
	resp, err := Request(d.URL, http.MethodGet, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(DownloadsDir, d.fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	d.bar.Describe("[cyan]下载中...[reset]")
	_, err = io.Copy(io.MultiWriter(file, d.bar), resp.Body)
	if err != nil {
		return err
	}
	err = d.bar.Finish()
	if err != nil {
		return err
	}
	return nil
}

func (d *Downloader) multiDownload() error {
	partSize := d.contentLength / d.maxConnections
	for partSize <= 1000000 && d.maxConnections > 1 {
		d.maxConnections--
		partSize = d.contentLength / d.maxConnections
	}

	// 创建部分文件的存放目录
	partDir := d.getPartDir()
	if err := os.Mkdir(partDir, 0777); err != nil {
		return err
	}
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(int(d.maxConnections))
	d.bar.Describe(fmt.Sprintf("[black]%d线程[cyan]同时下载中...[reset]", d.maxConnections))

	rangeStart := int64(0)
	for connectionNum := range d.maxConnections {
		go func(connectionNum, rangeStart int64) {
			defer wg.Done()

			rangeEnd := rangeStart + partSize
			// 最后一部分，总长度不能超过 ContentLength
			if connectionNum == d.maxConnections-1 {
				rangeEnd = d.contentLength
			}

			d.downloadPartial(rangeStart, rangeEnd, connectionNum)

		}(connectionNum, rangeStart)
	}
	wg.Wait()

	// 合并文件
	if err := d.merge(); err != nil {
		return err
	}

	if err := d.bar.Finish(); err != nil {
		return err
	}
	return nil
}

func (d *Downloader) downloadPartial(rangeStart, rangeEnd, connectionNum int64) {
	if rangeStart >= rangeEnd {
		return
	}
	resp, err := Request(d.URL, http.MethodGet, map[string]string{"Range": fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd-1)})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	partFile, err := os.OpenFile(d.getPartFilename(connectionNum), flags, 0666)
	if err != nil {
		panic(err)
	}
	defer partFile.Close()
	if _, err := io.Copy(partFile, resp.Body); err != nil {
		panic(err)
	}
}

func (d *Downloader) getPartDir() string {
	return filepath.Join(DownloadsDir, fmt.Sprintf("%s-parts", d.fileName))
}

// getPartFilename 构造部分文件的名字
func (d *Downloader) getPartFilename(partNum int64) string {
	return filepath.Join(d.getPartDir(), fmt.Sprintf("%s-%d", d.fileName, partNum))
}

func (d *Downloader) merge() error {
	destFile, err := os.OpenFile(filepath.Join(DownloadsDir, d.fileName), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for connectionNum := range d.maxConnections {
		partFileName := d.getPartFilename(connectionNum)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}
		if _, err := io.Copy(destFile, partFile); err != nil {
            return err
        }
		partFile.Close()
		os.Remove(partFileName)
	}

	return nil
}

func (d *Downloader) getFileName(header http.Header, url url.URL) {
	// 尝试从 Content-Disposition 头部获取文件名
	if disposition := header.Get("Content-Disposition"); disposition != "" {
		_, params, err := mime.ParseMediaType(disposition)
		if err == nil {
			if filename, ok := params["filename"]; ok {
				d.fileName = filename
				return
			}
		}
	}

	// 如果没有 Content-Disposition 头部，则从 URL 中获取文件名
	d.fileName = filepath.Base(url.Path)
}
