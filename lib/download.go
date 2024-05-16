package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

// 下载文件, 返回文件路径
func Download(url string, fileName string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "MCSCS-Golang/"+VERSION)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求%s失败, 状态码: %s", url, resp.Status)
	}

	// 进度条
	ansiStdout := ansi.NewAnsiStdout()
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
		}),
	)

	filePath := filepath.Join(DownloadsDir, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
