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
	"archive/zip"
	"bufio"
	"encoding/json"
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
	"strings"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

var EnableAria2c bool
var aria2cPath string

func InitAria2c() {
	var aria2cName string
	if runtime.GOOS == "windows" {
		aria2cName = "aria2c.exe"
	} else {
		aria2cName = "aria2c"
	}
	var err error
	aria2cPath, err = which("aria2c")
	if err != nil && os.IsNotExist(err) {
		if _, err := os.Stat(filepath.Join(aria2cDir, aria2cName)); os.IsNotExist(err) && runtime.GOOS == "windows" {
			aria2cPath = filepath.Join(aria2cDir, aria2cName)
		}
		if runtime.GOOS == "windows" {
			confirm, err := Confirm("aria2c 未安装, 是否下载(推荐)?")
			if err != nil {
				panic(err)
			}
			if confirm {
				err = downloadAria2c()
				if err != nil {
					panic(err)
				}
				EnableAria2c = true
				return
			}
		}
		EnableAria2c = false
	} else if err == nil {
		EnableAria2c = true
	}
	panic(err)
}

func NewDownloader(url url.URL) *Downloader {
	return &Downloader{URL: url}
}

type Downloader struct {
	URL      url.URL // 下载的 URL
	fileName string  // 文件名
}

func (d *Downloader) Download() (string, error) {
	resp, err := Request(d.URL, http.MethodGet, nil)
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
	ansiStdout := ansi.NewAnsiStdout()
	bar := progressbar.NewOptions64(
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
	f.WriteString(fmt.Sprintf(`%s
	referer=%s
	dir=%s
	out=%s`, d.URL.String(), d.URL.String(), DownloadsDir, d.fileName))
	err = f.Close()
	if err != nil {
		return err
	}
	cmd := exec.Command(aria2cPath)
	cmd.Args = append(cmd.Args,
		fmt.Sprintf("--input-file=%s", inputFilePath),
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
		d.URL.String(),
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

func Confirm(description string) (bool, error) {
	fmt.Printf("%s (y/n): ", description)
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	choice = strings.TrimSpace(strings.ToLower(choice))
	for {
		switch choice {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			continue
		}
	}
}

func downloadAria2c() error {
	if runtime.GOOS != "windows" {
		return errors.New("暂不支持非Windows系统")
	}
	aria2cPath = filepath.Join(aria2cDir, "aria2c.exe")
	resp, err := Request(url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "/repos/aria2/aria2/releases",
	}, http.MethodGet, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var data []struct {
		Assets []struct {
			Name        string `json:"name"`
			DownloadURL string `json:"browser_download_url"`
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	for _, asset := range data[0].Assets {
		if strings.Contains(asset.Name, "aria2c.exe") {
			url, err := url.Parse(asset.DownloadURL)
			if err != nil {
				return err
			}
			path, err := (&Downloader{URL: *url}).Download()
			if err != nil {
				return err
			}
			r, err := zip.OpenReader(path)
			if err != nil {
				return err
			}
			// 解压文件
			for _, f := range r.File {
				rc, err := f.Open()
				if err != nil {
					return err
				}
				defer rc.Close()
				path := filepath.Join(aria2cDir, f.Name)
				if f.FileInfo().IsDir() {
					os.MkdirAll(path, f.Mode())
				} else {
					dir := filepath.Dir(path)
					if err := os.MkdirAll(dir, 0755); err != nil {
						return err
					}
					dst, err := os.Create(path)
					if err != nil {
						return err
					}
					defer dst.Close()
					if _, err := io.Copy(dst, rc); err != nil {
						return err
					}
				}
			}
			return nil
		}
	}
	return nil
}

func which(command string) (string, error) {
	path, err := exec.LookPath(command)
	if err != nil {
		return "", err
	}
	return path, nil
}
