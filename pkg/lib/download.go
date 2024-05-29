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
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

type Downloader struct {
	URL           url.URL                  // 下载的 URL
	FileName      string                   // 保存的文件名
	Multi         int                      // 是否多线程下载: 0 自动选择(默认), 1 单线程下载, 2 多线程下载
	bar           *progressbar.ProgressBar // 进度条
	maxConnetions int                      // 最大连接数
	contentLength int                      // 内容长度
}

func (d Downloader) Download() (string, error) {
	resp, err := Request(d.URL, http.MethodHead, nil)
	if err != nil {
		return "", err
	}

	path := filepath.Join(downloadsDir, d.FileName)
	if resp.StatusCode != http.StatusOK {
		return path, fmt.Errorf(resp.Status)
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
	d.maxConnetions = configs.MaxConnetions
	d.contentLength = int(resp.ContentLength)
	switch d.Multi {
	case 1:
		return path, d.singleDownload()
	case 2:
		if resp.Header.Get("Accept-Ranges") != "bytes" {
			log.Warn().Msgf("下载 %s 时不支持多线程下载", d.URL.String())
			return path, d.singleDownload()
		}
		d.multiDownload()
	default:
		if resp.Header.Get("Accept-Ranges") == "bytes" && resp.ContentLength > 1024*1024 {
			return path, d.multiDownload()
		}
		return path, d.singleDownload()
	}
	return path, nil
}

func (d Downloader) singleDownload() error {
	resp, err := Request(d.URL, http.MethodGet, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(downloadsDir, d.FileName))
	if err != nil {
		return err
	}
	defer file.Close()

	d.bar.Describe("[cyan]下载中...[reset]")
	_, err = io.Copy(io.MultiWriter(file, d.bar), resp.Body)
	if err != nil {
		return err
	}
	d.bar.Finish()
	return nil
}

func (d Downloader) multiDownload() error {
	partSize := d.contentLength / d.maxConnetions

	// 创建部分文件的存放目录
	partDir := d.getPartDir()
	os.Mkdir(partDir, 0777)
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(d.maxConnetions)
	d.bar.Describe(fmt.Sprintf("[black]%d线程[cyan]同时下载中...[reset]", d.maxConnetions))

	rangeStart := 0
	for connectionNum := range d.maxConnetions {
		go func(connectionNum, rangeStart int) {
			defer wg.Done()

			rangeEnd := rangeStart + partSize
			// 最后一部分，总长度不能超过 ContentLength
			if connectionNum == d.maxConnetions-1 {
				rangeEnd = d.contentLength
			}

			d.downloadPartial(rangeStart, rangeEnd, connectionNum)

		}(connectionNum, rangeStart)
	}

	wg.Wait()

	// 合并文件
	d.merge()

	d.bar.Finish()

	return nil
}

func (d Downloader) downloadPartial(rangeStart, rangeEnd, connectionNum int) {
	if rangeStart >= rangeEnd {
		return
	}
	resp, err := Request(d.URL, http.MethodGet, map[string]string{"Range": fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd-1)})
	if err != nil {
		log.Error().Err(err).Msgf("下载部分文件 %d 失败", connectionNum)
		return
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	partFile, err := os.OpenFile(d.getPartFilename(connectionNum), flags, 0666)
	if err != nil {
		log.Error().Err(err).Msgf("打开部分文件 %d 失败", connectionNum)
		return
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(partFile, d.bar), resp.Body, buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Error().Err(err).Msgf("写入部分文件 %d 失败", connectionNum)
		return
	}
}

func (d Downloader) getPartDir() string {
	return filepath.Join(downloadsDir, fmt.Sprintf("%s-parts", d.FileName))
}

// getPartFilename 构造部分文件的名字
func (d Downloader) getPartFilename(partNum int) string {
	return filepath.Join(d.getPartDir(), fmt.Sprintf("%s-%d", d.FileName, partNum))
}

func (d Downloader) merge() error {
	destFile, err := os.OpenFile(filepath.Join(downloadsDir, d.FileName), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for connectionNum := range d.maxConnetions {
		partFileName := d.getPartFilename(connectionNum)
		partFile, err := os.Open(partFileName)
		if err != nil {
			return err
		}
		io.Copy(destFile, partFile)
		partFile.Close()
		os.Remove(partFileName)
	}

	return nil
}
