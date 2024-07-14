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
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Arama0517/MCST/internal/build"
	"github.com/Arama0517/MCST/internal/configs"
	"github.com/Arama0517/MCST/internal/requests"
	"github.com/apex/log"
	"github.com/schollz/progressbar/v3"
	"github.com/siku2/arigo"
)

type Downloader struct {
	URL      string
	FileName string
	bar      *progressbar.ProgressBar
}

func NewDownloader(url string) *Downloader {
	return &Downloader{
		URL: url,
		bar: progressbar.NewOptions64(
			100,
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
				_, err := fmt.Fprint(os.Stderr, "\n")
				if err != nil {
					return
				}
			})),
	}
}

func (d *Downloader) Download() (string, error) {
	req, err := requests.NewRequest(http.MethodGet, d.URL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()
	d.FileName = resp.Header.Get("Content-Disposition")
	if d.FileName == "" {
		d.FileName = filepath.Base(resp.Request.URL.Path)
	}
	filePath := filepath.Join(configs.DownloadsDir, d.FileName)
	switch configs.Configs.Settings.Downloader {
	default:
	case 1:
		return d.aria2Download()
	case 2:
		if runtime.GOOS == "windows" {
			return d.idmDownload()
		}
		log.Warn("非 Windows 系统暂不支持 IDM 下载方式")
	}
	return d.defaultDownload(filePath, resp)
}

func (d *Downloader) defaultDownload(filePath string, resp *http.Response) (string, error) {
	d.bar.ChangeMax64(resp.ContentLength)
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

func (d *Downloader) aria2Download() (string, error) {
	aria2Name := "aria2c"
	if runtime.GOOS == "windows" {
		aria2Name += ".exe"
	}
	cmd := exec.Command(aria2Name)
	cmd.Args = append(cmd.Args, configs.Configs.Settings.Aria2.Options...)
	cmd.Args = append(cmd.Args,
		"--dir="+configs.DownloadsDir,
		fmt.Sprintf("--user-agent=MCST/%s", build.Version.GitVersion),
		"--allow-overwrite=true",
		"--auto-file-renaming=false",
		fmt.Sprintf("--retry-wait=%d", configs.Configs.Settings.Aria2.RetryWait),
		fmt.Sprintf("--split=%d", configs.Configs.Settings.Aria2.Split),
		fmt.Sprintf("--max-connection-per-server=%d", configs.Configs.Settings.Aria2.MaxConnectionPerServer),
		fmt.Sprintf("--min-split-size=%s", configs.Configs.Settings.Aria2.MinSplitSize),
		"--enable-rpc",
		"--rpc-listen-all",
		"--rpc-listen-port=6800",
		"--rpc-secret=MCST",
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
	if err := cmd.Start(); err != nil {
		return "", err
	}

	time.Sleep(500 * time.Millisecond)
	client, err := arigo.Dial("ws://127.0.0.1:6800", "MCST")
	if err != nil {
		return "", err
	}
	gid, err := client.AddURI(arigo.URIs(d.URL), nil)
	if err != nil {
		return "", err
	}

	time.Sleep(500 * time.Millisecond)
	status, err := gid.TellStatus("totalLength")
	if err != nil {
		return "", err
	}
	d.bar.ChangeMax64(int64(status.TotalLength))

	for {
		status, err = gid.TellStatus("status", "completedLength", "connections")
		if err != nil {
			return "", err
		}
		if status.Status == "complete" {
			files, err := gid.GetFiles()
			if err != nil {
				return "", err
			}
			return files[0].Path, nil
		}
		if err = d.bar.Set64(int64(status.CompletedLength)); err != nil {
			return "", err
		}
		time.Sleep(50 * time.Millisecond)
	}
}
