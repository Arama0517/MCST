# Minecraft Server Tool

- Powered by **[Golang](https://go.dev/)**

## 安装

### 使用Shell脚本安装
```shell
curl -sSfL https://raw.githubusercontent.com/Arama-Vanarana/MCServerTool/main/scripts/install.sh | sudo sh -s -- -b /usr/local/bin
```

### Go install
> [Go 1.22+](https://go.dev/doc/install)
```shell
go install github.com/Arama0517/MCST@latest
```

### deb, rpm, zst 和 apk 包
- 从[Releases](https://github.com/Arama-Vanarana/MCServerTool/releases/latest)下载, 并使用合适的工具安装他们

### homebrew
```shell
brew install Arama0517/tap/MCST
```

### scoop
```shell
scoop bucket add MCST https://github.com/Arama0517/scoop-bucket.git
scoop install MCST
```

### 从源代码开始构建
这里你有2个选择:
1. 如果你想为该项目做出贡献, 请按照我们的[贡献指南](./CONTRIBUTING.md)中的步骤进行操作
2. 如果你处于某种原因只是想从源代码构建, 请安装以下步骤操作

**克隆**:
```shell
git clone https://github.com/Arama0517/MCServerTool.git --depth 1
```
**下载依赖**:
```shell
go mod tidy
```
**构建**:
```shell
go build -o MCST .
```
**检查他是否有效**:
```shell
./MCST --version
```