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
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

var concurrency int
var ansiStdout = ansi.NewAnsiStdout()

func Download(URL url.URL, filename string) (string, error) {
	resp, err := Request(URL, http.MethodHead, nil)
	if err != nil {
		return "", err
	}
	configs, err := LoadConfigs()
	if err != nil {
		return "", err
	}
	concurrency = configs.Concurrency
	path := filepath.Join(DownloadsDir, filename)
	if resp.Header.Get("Accept-Ranges") == "bytes" && resp.ContentLength > 1024*1024 {
		return path, MultiDownload(URL, filename, int(resp.ContentLength))
	}
	return path, SingleDownload(URL, filename)
}

func SingleDownload(URL url.URL, filename string) error {
	resp, err := Request(URL, http.MethodGet, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath.Join(DownloadsDir, filename))
	if err != nil {
		return err
	}
	defer file.Close()

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
		}))
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func MultiDownload(URL url.URL, filename string, contentLen int) error {
	partSize := contentLen / concurrency

	// 创建部分文件的存放目录
	partDir := getPartDir(filename)
	os.Mkdir(partDir, 0777)
	defer os.RemoveAll(partDir)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	rangeStart := 0
	bar := progressbar.NewOptions(
		contentLen,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription(fmt.Sprintf("[black]%d线程[cyan]同时下载中...[reset]", concurrency)),
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

	for i := 0; i < concurrency; i++ {
		// 并发请求
		go func(i, rangeStart int) {
			defer wg.Done()

			rangeEnd := rangeStart + partSize
			// 最后一部分，总长度不能超过 ContentLength
			if i == concurrency-1 {
				rangeEnd = contentLen
			}

			downloadPartial(URL, filename, rangeStart, rangeEnd, i, bar)

		}(i, rangeStart)

		rangeStart += partSize + 1
	}

	wg.Wait()

	// 合并文件
	merge(filename)

	return nil
}

func downloadPartial(URL url.URL, filename string, rangeStart, rangeEnd, i int, bar *progressbar.ProgressBar) {
	if rangeStart >= rangeEnd {
		return
	}
	resp, err := Request(URL, http.MethodGet, map[string]string{"Range": fmt.Sprintf("bytes=%d-%d", rangeStart, rangeEnd-1)})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	flags := os.O_CREATE | os.O_WRONLY
	partFile, err := os.OpenFile(getPartFilename(filename, i), flags, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer partFile.Close()

	buf := make([]byte, 32*1024)
	_, err = io.CopyBuffer(io.MultiWriter(partFile, bar), resp.Body, buf)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Fatal(err)
	}
}

// getPartDir 构造部分文件的存放目录
func getPartDir(filename string) string {
	return filepath.Join(DownloadsDir, fmt.Sprintf("%s-parts", filename))
}

// getPartFilename 构造部分文件的名字
func getPartFilename(filename string, partNum int) string {
	return filepath.Join(getPartDir(filename), fmt.Sprintf("%s-%d", filename, partNum))
}

func merge(filename string) error {
	destFile, err := os.OpenFile(filepath.Join(DownloadsDir, filename), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for i := 0; i < concurrency; i++ {
		partFileName := getPartFilename(filename, i)
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
