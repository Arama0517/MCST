# Minecraft Server Tool

- Powered by **[Golang](https://go.dev/)**

## 安装

### Windows

1. 在[Releases](https://github.com/Arama-Vanarana/MCServerTool/releases/latest)下载最新版
2. 解压到任意目录

### Linux

#### 通过下载二进制安装

> 请在curl和wget中选一个下载方式

##### curl

```shell
# 如果你在root用户运行命令请去掉sudo
curl -sSfL https://raw.githubusercontent.com/Arama-Vanarana/MCServerTool/main/scripts/install.sh | sudo sh -s -- -b /usr/local/bin
```

##### wget

```shell
# 如果你在root用户运行命令请去掉sudo
wget -O- -nv https://raw.githubusercontent.com/Arama-Vanarana/MCServerTool/main/scripts/install.sh | sudo sh -s -- -b /usr/local/bin
```

#### 从源代码编译

1. 安装[Git](https://www.git-scm.com/downloads)和[Go](https://go.dev/dl)
2. 运行命令

```shell
# 1. 克隆仓库并进入
git clone https://github.com/Arama-Vanarana/MCServerTool.git
cd MCServerTool
# 2. 编译程序
go build -s -w -o MCST ./cmd/MCST/main.go 
# 3. 安装到 /usr/local/bin
## 如果你在root用户运行命令请去掉sudo
sudo cp MCST /usr/local/bin
# 4. (可选) 删除仓库源代码
cd ..
rm -r MCServerTool
```

## 卸载

### Windows

- 删除解压程序的文件夹

### Linux

- 运行命令

```shell
# 如果你在root用户运行命令请去掉sudo
sudo rm /usr/local/bin/MCST
```