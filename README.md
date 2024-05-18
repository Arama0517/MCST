<!-- 
MCSCS can be used to easily create, launch, and configure a Minecraft server.
Copyright (C) 2024 Arama

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. 
-->


# MCSCS

-   Powered by*[Golang](https://go.dev/)**

# 安装

## Windows

1. 在[Releases](../../releases/latest)下载最新版(windows-X.X.X-amd64.zip)
2. 解压到任意目录
3. 运行`MCSCS.exe`

## Linux

1. 安装[make](https://www.gnu.org/software/make), [go](https://go.dev/dl)

    > 推荐使用 [GVM](https://github.com/moovweb/gvm) 进行下载 go

2. 运行以下命令

```bash
# 克隆仓库
git clone https://github.com/Arama-Vanarana/MCSCS-Go.git
cd MCSCS-Go
make install
cd ..
rm -rf MCSCS-Go 
# 运行MCSCS
mcscs
```
