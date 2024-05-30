package pages

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Arama-Vanarana/MCServerTool/pkg/lib"
	"github.com/urfave/cli/v2"
)

func confirm(description string) (bool, error) {
	fmt.Printf("%s (y/n): ", description)
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
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

func Create(c *cli.Context) error {
	// 配置
	configs, err := lib.LoadConfigs()
	if err != nil {
		return err
	}
	var config lib.ServerConfig
	config.Name = c.String("name")
	for name := range configs.Servers {
		if name == config.Name {
			return errors.New("服务器已存在")
		}
	}
	config.Java.Xms, err = ToBytes(c.String("Xms"))
	if err != nil {
		return err
	}
	config.Java.Xmx, err = ToBytes(c.String("Xmx"))
	if err != nil {
		return err
	}
	config.Java.Encoding = "UTF-8"
	if c.Bool("gbk") {
		config.Java.Encoding = "GBK"
	}
	config.Java.Path = c.Path("java")
	config.Java.Args = strings.Split(c.String("jvm_args"), " ")
	config.ServerArgs = strings.Split(c.String("server_args"), " ")
	coreIndex := c.Int("core")
	if coreIndex < 0 || coreIndex >= len(configs.Cores) {
		return errors.New("核心不存在")
	}
	coreInfo := configs.Cores[coreIndex]
	configs.Servers[config.Name] = config
	// EULA
	serverPath := filepath.Join(lib.ServersDir, config.Name)
	if err := os.MkdirAll(serverPath, 0755); err != nil {
		return err
	}
	choice, err := confirm("你是否同意EULA协议<https://aka.ms/MinecraftEULA/>?")
	if err != nil {
		return err
	}
	if choice {
		file, err := os.Create(filepath.Join(serverPath, "eula.txt"))
		if err != nil {
			return err
		}
		defer file.Close()
		data := fmt.Sprintf(`# Create By Minecraft Server Tool
# By changing the setting below to TRUE you are indicating your agreement to Minecraft EULA(<https://aka.ms/MinecraftEULA/>).
# %s
eula=true`, time.Now().Format("Mon Jan 02 15:04:05 MST 2006"))
		file.WriteString(data)
	} else {
		return errors.New("你必须同意EULA协议<https://aka.ms/MinecraftEULA/>才能创建服务器")
	}
	// 保存
	if err := configs.Save(); err != nil {
		return err
	}
	// 复制
	srcFile, err := os.Open(coreInfo.FilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(filepath.Join(serverPath, "server.jar"))
	if err != nil {
		return err
	}
	defer dstFile.Close()
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	if err := os.Chmod(filepath.Join(serverPath, "server.jar"), 0755); err != nil {
		return err
	}
	return nil
}

var ErrorNotAUnit = errors.New("这不是一个有效的单位")

func ToBytes(byteStr string) (uint64, error) {
	var num, unit string

	// 分离数字部分和单位部分
	for _, char := range byteStr {
		if char >= '0' && char <= '9' {
			num += string(char)
		} else if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
			unit += strings.ToUpper(string(char))
		} else {
			return 0, nil
		}
	}

	var multiplier uint64
	switch unit {
	case "T", "TB", "TIB":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "G", "GB", "GIB":
		multiplier = 1024 * 1024 * 1024
	case "M", "MB", "MIB":
		multiplier = 1024 * 1024
	case "K", "KB", "KIB":
		multiplier = 1024
	case "B", "BYTES", "":
		multiplier = 1
	default:
		return 0, ErrorNotAUnit
	}

	parsedNum, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}

	return parsedNum * multiplier, nil
}
