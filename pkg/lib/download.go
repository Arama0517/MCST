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
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

var (
	EnableAria2c bool
	aria2cPath   string
)

func initDownloader() error {
	aria2cName := "aria2c"
	if runtime.GOOS == "windows" {
		aria2cName = "aria2c.exe"
	}
	var err error
	aria2cPath, err = exec.LookPath(aria2cName)
	switch err {
	case nil:
		EnableAria2c = true
	case os.ErrNotExist:
		EnableAria2c = false
	default:
		return err
	}
	return nil
}

func NewDownloader(url url.URL) *Downloader {
	return &Downloader{URL: url}
}

type Downloader struct {
	URL      url.URL // 下载的 URL
	fileName string  // 文件名
}

func (d *Downloader) Download() (string, error) {
	resp, err := Request(d.URL, http.MethodGet, nil, nil)
	if err != nil {
		return "", err
	}

	d.getFileName(resp.Header, d.URL)
	path := filepath.Join(DownloadsDir, d.fileName)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	if EnableAria2c {
		return path, d.aria2cDownload()
	}
	// 单线程
	ansiStderr := ansi.NewAnsiStderr()
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionSetWriter(ansiStderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(ansiStderr, "\n")
		}))
	file, err := os.Create(filepath.Join(DownloadsDir, d.fileName))
	if err != nil {
		return "", err
	}
	defer file.Close()

	bar.Describe("[cyan]下载中...[reset]")
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		return "", err
	}
	err = bar.Finish()
	if err != nil {
		return "", err
	}
	return path, nil
}

func (d *Downloader) aria2cDownload() error {
	configs, err := LoadConfigs()
	if err != nil {
		return err
	}
	inputFilePath := filepath.Join(DownloadsDir, fmt.Sprintf("%s.txt", d.fileName))
	f, err := os.Create(inputFilePath)
	if err != nil {
		return err
	}
	// f.WriteString(fmt.Sprintf(`%s
	// dir=%s
	// out=%s`, d.URL.String(), DownloadsDir, d.fileName))
	f.WriteString(d.URL.String())
	if d.URL.Host != "sourceforge.net" {
		f.WriteString("    referer=" + filepath.Dir(d.URL.Path))
	}
	f.WriteString("    dir=" + DownloadsDir)
	f.WriteString("    out=" + d.fileName)
	err = f.Close()
	if err != nil {
		return err
	}
	cmd := exec.Command(aria2cPath)
	cmd.Args = append(cmd.Args,
		fmt.Sprintf("--input-file='%s'", inputFilePath),
		fmt.Sprintf("--user-agent='MCST/%s'", Version),
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		fmt.Sprintf("--retry-wait=%d", configs.Aria2c.RetryWait),
		fmt.Sprintf("--max-connection-per-server=%d", configs.Aria2c.MaxConnectionPerServer),
		fmt.Sprintf("--min-split-size=%s", configs.Aria2c.MinSplitSize),
		"--console-log-level=warn",
		"--no-conf=true",
		"--follow-metalink=true",
		"--metalink-preferred-protocol=https",
		"--min-tls-version=TLSv1.2",
		fmt.Sprintf("--stop-with-process=%d", os.Getpid()),
		"--continue",
		"--summary-interval=0",
		"--auto-save-interval=1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		if err := os.Remove(inputFilePath); err != nil {
			return err
		}
		return err
	}
	if err := os.Remove(inputFilePath); err != nil {
		return err
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
