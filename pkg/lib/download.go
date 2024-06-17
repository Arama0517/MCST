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

import "C"

import (
	"context"
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
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/ybbus/jsonrpc/v3"
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
	// 检测是否已下载
	resp, err := Request(d.URL, http.MethodGet, nil, nil)
	if err != nil {
		return "", err
	}
	path := filepath.Join(DownloadsDir, d.FileName)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	d.getFileName(resp.Header)

	// Aria2 多线程下载
	if Configs.Aria2c.Enabled {
		aria2cName := "aria2c"
		if runtime.GOOS == "windows" {
			aria2cName = "aria2c.exe"
		}
		var err error
		if d.aria2Path, err = exec.LookPath(aria2cName); err != nil && !errors.Is(err, exec.ErrNotFound) {
			return "", err
		} else if err == nil {
			return path, d.aria2cDownload()
		}
	}
	// 单线程
	bar := progressbar.NewOptions64(
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

//nolint:tagliatelle // Aria2 jsonrpc返回
type downloadStatus struct {
	Status          string `json:"status"`
	TotalLength     string `json:"totalLength"`
	CompletedLength string `json:"completedLength"`
	DownloadSpeed   string `json:"downloadSpeed"`
	Connections     string `json:"connections"`
	ErrorMessage    string `json:"errorMessage"`
}

func (d *Downloader) aria2cDownload() error {
	cmd := exec.Command(d.aria2Path)
	cmd.Args = append(cmd.Args,
		"--dir="+DownloadsDir,
		fmt.Sprintf("--user-agent=MCST/%s", version.GitVersion),
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		fmt.Sprintf("--retry-wait=%d", Configs.Aria2c.RetryWait),
		fmt.Sprintf("--split=%d", Configs.Aria2c.Split),
		fmt.Sprintf("--max-connection-per-server=%d", Configs.Aria2c.MaxConnectionPerServer),
		fmt.Sprintf("--min-split-size=%s", Configs.Aria2c.MinSplitSize),
		"--enable-rpc",
		// "--console-log-level=error",
		"--quiet",
		"--no-conf=true",
		"--follow-metalink=true",
		"--metalink-preferred-protocol=https",
		"--min-tls-version=TLSv1.2",
		fmt.Sprintf("--stop-with-process=%d", os.Getpid()),
		"--continue",
		"--summary-interval=0",
		"--auto-save-interval=1",
	)
	cmd.Args = append(cmd.Args, Configs.Aria2c.Option...)
	if err := cmd.Start(); err != nil {
		return err
	}
	defer func(Process *os.Process) {
		_ = Process.Kill()
	}(cmd.Process)
	time.Sleep(500 * time.Millisecond)
	ctx := context.Background()
	client := jsonrpc.NewClient("http://127.0.0.1:6800/jsonrpc")
	var gid string
	if err := client.CallFor(ctx, &gid, "aria2.addUri", []any{[]string{d.URL.String()}}); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	var status downloadStatus
	if err := client.CallFor(ctx, &status, "aria2.tellStatus", []any{gid, []string{
		"status",
		"totalLength",
		"completedLength",
		"downloadSpeed",
		"connections",
		"errorMessage",
	}}); err != nil {
		return err
	}
	totalLength, err := strconv.ParseInt(status.TotalLength, 10, 64)
	if err != nil {
		return err
	}
	completedLength, err := strconv.ParseInt(status.CompletedLength, 10, 64)
	if err != nil {
		return err
	}
	bar := progressbar.NewOptions64(
		totalLength,
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
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			_, err := fmt.Fprint(os.Stderr, "\n")
			if err != nil {
				return
			}
		}))
	for {
		if err := client.CallFor(ctx, &status, "aria2.tellStatus", []any{gid, []string{
			"status",
			"totalLength",
			"completedLength",
			"downloadSpeed",
			"connections",
			"errorMessage",
		}}); err != nil {
			return err
		}
		if status.Status == "complete" {
			return bar.Finish()
		} else if status.Status == "error" {
			return errors.New(status.ErrorMessage)
		}
		if err := bar.Set64(completedLength); err != nil {
			return err
		}
		bar.Describe(fmt.Sprintf("[cyan]%s服务器同时下载中...[reset]", status.Connections))
		time.Sleep(500 * time.Millisecond)
	}
}

func (d *Downloader) getFileName(header http.Header) {
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
	d.FileName = filepath.Base(d.URL.Path)
}
