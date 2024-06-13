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

package lib

import (
	"errors"
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

	"github.com/schollz/progressbar/v3"
)

func NewDownloader(url url.URL) *Downloader {
	return &Downloader{URL: url}
}

type Downloader struct {
	URL       url.URL // 下载的 URL
	FileName  string  // 文件名
	aria2Path string
}

func (d *Downloader) Download() (string, error) {
	resp, err := Request(d.URL, http.MethodGet, nil, nil)
	if err != nil {
		return "", err
	}
	d.getFileName(resp.Header, d.URL)
	path := filepath.Join(DownloadsDir, d.FileName)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	if configs, err := LoadConfigs(); err == nil && configs.Aria2c.Enabled {
		aria2cName := "aria2c"
		if runtime.GOOS == "windows" {
			aria2cName = "aria2c.exe"
		}
		switch d.aria2Path, err = exec.LookPath(aria2cName); {
		case errors.Is(err, nil), errors.Is(err, exec.ErrNotFound):
			break
		default:
			return "", err
		}
		if _, err = os.Stat(d.aria2Path); err == nil {
			if err := resp.Body.Close(); err != nil {
				return "", err
			}
			return path, d.aria2cDownload()
		}
	}
	// 单线程
	bar := progressbar.NewOptions64(
		resp.ContentLength,
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
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			_, err := fmt.Fprint(os.Stderr, "\n")
			if err != nil {
				return
			}
		}))
	file, err := os.Create(filepath.Join(DownloadsDir, d.FileName))
	if err != nil {
		return "", err
	}
	bar.Describe("[cyan]下载中...[reset]")
	if _, err := io.Copy(io.MultiWriter(file, bar), resp.Body); err != nil {
		return "", err
	}
	if err := resp.Body.Close(); err != nil {
		return "", err
	}
	if err := bar.Finish(); err != nil {
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	return path, nil
}

func (d *Downloader) aria2cDownload() error {
	inputFilePath := filepath.Join(DownloadsDir, fmt.Sprintf("%s.txt", d.FileName))
	f, err := os.Create(inputFilePath)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(d.URL.String() + "\n"); err != nil {
		return err
	}
	if d.URL.Host != "sourceforge.net" {
		if _, err := f.WriteString(fmt.Sprintf("\treferer=%s://%s%s\n", d.URL.Scheme, d.URL.Host, filepath.Dir(d.URL.Path))); err != nil {
			return err
		}
	}
	if _, err := f.WriteString(fmt.Sprintf("\tdir=%s\n", DownloadsDir)); err != nil {
		return err
	}
	if _, err := f.WriteString(fmt.Sprintf("\tout=%s\n", d.FileName)); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(inputFilePath)
	configs, err := LoadConfigs()
	if err != nil {
		return err
	}
	cmd := exec.Command(d.aria2Path)
	cmd.Args = append(cmd.Args,
		fmt.Sprintf("--input-file=%s", inputFilePath),
		fmt.Sprintf("--user-agent=MCST/%s", version),
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		fmt.Sprintf("--retry-wait=%d", configs.Aria2c.RetryWait),
		fmt.Sprintf("--split=%d", configs.Aria2c.Split),
		fmt.Sprintf("--max-connection-per-server==%d", configs.Aria2c.MaxConnectionPerServer),
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
	cmd.Args = append(cmd.Args, configs.Aria2c.Option...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (d *Downloader) getFileName(header http.Header, url url.URL) {
	// 尝试从 Content-Disposition 头部获取文件名
	if disposition := header.Get("Content-Disposition"); disposition != "" {
		_, params, err := mime.ParseMediaType(disposition)
		if err == nil {
			if filename, ok := params["filename"]; ok {
				d.FileName = filename
				return
			}
		}
	}
	// 如果没有 Content-Disposition 头部，则从 URL 中获取文件名
	d.FileName = filepath.Base(url.Path)
}
