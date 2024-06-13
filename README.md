# Minecraft Server Tool

- Powered by **[Golang](https://go.dev/)**

## 安装

### Windows

1. 在[Releases](https://github.com/Arama-Vanarana/MCServerTool/releases/latest)下载最新版
2. 解压到任意目录

### Linux

#### 使用Shell脚本安装
```shell
# 如果你在root用户运行命令请去掉sudo
curl -sSfL https://raw.githubusercontent.com/Arama-Vanarana/MCServerTool/main/scripts/install.sh | sh -s -- -b /usr/bin
```

#### Go install (Go 1.22)
```shell
go install github.com/Arama0517/MCServerTool@latest
```

#### deb, rpm, zst 和 apk 包
* 从[Releases](https://github.com/Arama-Vanarana/MCServerTool/releases/latest)下载, 并使用合适的工具安装他们

## 卸载

### Windows

- 删除解压程序的文件夹

### Linux

- 运行命令

```shell
# 如果你在root用户运行命令请去掉sudo
sudo rm /usr/local/bin/MCST
```