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
	"github.com/schollz/progressbar/v3"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var (
	EnableAria2c bool
	aria2cPath   string
)

func initDownloader() error {
	configs, err := LoadConfigs()
	if err != nil {
		return err
	}
	aria2cPath = configs.Aria2c.Path
	if aria2cPath == "auto" {
		aria2cName := "aria2c"
		if runtime.GOOS == "windows" {
			aria2cName = "aria2c.exe"
		}
		aria2cPath, err = exec.LookPath(aria2cName)
		switch {
		case err == nil:
			EnableAria2c = true
		case errors.Is(err, exec.ErrNotFound):
			EnableAria2c = false
		default:
			return err
		}
	}
	return nil
}

func NewDownloader(url url.URL) *Downloader {
	return &Downloader{URL: url, Stdout: os.Stdout, Stderr: os.Stderr}
}

type Downloader struct {
	URL            url.URL   // 下载的 URL
	Stdout, Stderr io.Writer // 输出
	fileName       string    // 文件名
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
		progressbar.OptionSetWriter(d.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			_, err := fmt.Fprint(d.Stderr, "\n")
			if err != nil {
				return
			}
		}))
	file, err := os.Create(filepath.Join(DownloadsDir, d.fileName))
	if err != nil {
		return "", err
	}
	bar.Describe("[cyan]下载中...[reset]")
	if _, err := io.Copy(io.MultiWriter(file, bar), resp.Body); err != nil {
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
	configs, err := LoadConfigs()
	if err != nil {
		return err
	}
	inputFilePath := filepath.Join(DownloadsDir, fmt.Sprintf("%s.txt", d.fileName))
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
	if _, err := f.WriteString(fmt.Sprintf("\tout=%s\n", d.fileName)); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	cmd := exec.Command(aria2cPath)
	cmd.Args = append(cmd.Args,
		fmt.Sprintf("--input-file=%s", inputFilePath),
		fmt.Sprintf("--user-agent=%s", userAgent),
		fmt.Sprintf("--stop-with-process=%d", os.Getpid()),
	)
	cmd.Args = append(cmd.Args, configs.Aria2c.Args...)
	cmd.Stdout = d.Stdout
	cmd.Stderr = d.Stderr
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
