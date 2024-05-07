package lib

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var Logger = logrus.New()
var aria2Process *os.Process

func init() {
	initLogs()
	initDataDirs()
	initAria2c()
}

// 初始化日志
func initLogs() {
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetReportCaller(true)
	Logger.SetFormatter(joonix.NewFormatter())
	logFile, err := os.OpenFile(filepath.Join(GetLogsDir(), "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	Logger.SetOutput(logFile)
	Logger.Info("Hello world!")
	Logger.Info("本程序遵循GPLv3协议开源")
	Logger.Info("作者: Arama 3584075812@qq.com")
}

// 初始化aria2c
func initAria2c() {
	aria2cDir := GetAria2cDir()
	configPath := filepath.Join(aria2cDir, "aria2c.conf")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, err := os.Create(configPath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		resp, err := http.Get("https://github.com/Arama-Vanarana/MCSCS-Golang/releases/download/aria2/aria2c.conf")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			panic(err)
		}
	}

	aria2cPath, err := exec.LookPath("aria2c")
	if err != nil {
		var executable string
		if runtime.GOOS == "windows" {
			executable = "aria2c.exe"
		} else {
			executable = "aria2c"
		}
		aria2cPath = filepath.Join(aria2cDir, executable)
		if _, err := os.Stat(aria2cPath); os.IsNotExist(err) {
			if runtime.GOOS == "windows" {
				client := &http.Client{}
				request, err := http.NewRequest("GET", "https://github.com/Arama-Vanarana/MCSCS-Golang/releases/download/aria2/aria2c.conf", nil)
				if err != nil {
					panic(err)
				}
				request.Header.Set("User-Agent", "MCSCS/"+VERSION)
				response, err := client.Do(request)
				if err != nil {
					panic(err)
				}
				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)
				if err != nil {
					panic(err)
				}
				var data map[string]interface{}
				err = json.Unmarshal(body, &data)
				if err != nil {
					panic(err)
				}
				var url string = "unknwon"
				if value, ok := data["key"].([]interface{}); ok {
					for _, item := range value {
						if obj, ok := item.(map[string]interface{}); ok {
							if value, ok := obj["browser_download_url"]; ok {
								url = value.(string)
							}
						}
					}
				}
				if url == "unknwon" {
					panic("获取aria2c下载链接失败, 请查看你的网络是否存在问题后重新开启本程序.")
				}
				file, err := os.Create(filepath.Join(aria2cDir, "aria2c.exe"))
				if err != nil {
					panic(err)
				}
				defer file.Close()
				resp, err := http.Get(url)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()
				_, err = io.Copy(file, resp.Body)
				if err != nil {
					panic(err)
				}
			} else if runtime.GOOS == "linux" {
				cfg, err := ini.Load("/etc/os-release")
				if err != nil {
					panic(err)
				}
				section, err := cfg.GetSection("")
				if err != nil {
					panic(err)
				}
				var cmd *exec.Cmd
				switch section.Key("ID").String() {
				case "debian", "ubuntu":
					cmd = exec.Command("sudo", "apt", "install", "aria2", "-y")
				case "centos", "fedora", "rhel":
					cmd = exec.Command("sudo", "yum", "install", "aria2", "-y")
				case "arch", "manjaro":
					cmd = exec.Command("sudo", "pacman", "-S", "--noconfirm", "aria2")
				default:
					panic("不支持自动下载的系统, 请按照<https://github.com/aria2/aria2/>安装aria2c")
				}
				cmd.Stdout = os.Stdout
				err = cmd.Run()
				if err != nil {
					panic(err)
				}
			} else {
				panic("不支持自动下载的系统, 请按照<https://github.com/aria2/aria2/>安装aria2c")
			}
		}
	}
	cmd := exec.Command(aria2cPath, "--conf-path="+configPath, "--dir="+aria2cDir, "--log="+filepath.Join(GetLogsDir(), "aria2c.log"), "--enable-rpc=true", "--rpc-listen-port=6800", "--quiet=true")
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	aria2Process = cmd.Process
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		KillAria2c()
		os.Exit(0)
	}()
}

func initDataDirs() {
	Logger.Info("数据根目录: ", GetDataDir())
	Logger.Info("配置存放目录: ", GetConfigsDir())
	Logger.Info("服务器存放目录: ", GetServersDir())
	Logger.Info("Aria2c配置/可执行程序存放目录: ", GetAria2cDir())
	Logger.Info("日志存放目录: ", GetLogsDir())
}
