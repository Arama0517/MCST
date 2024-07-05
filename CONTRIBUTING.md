# 贡献

参与本项目即表示你同意遵守我们的
[行为准则](https://github.com/Arama0517/MCST/blob/main/CODE_OF_CONDUCT.md)

## 设置你的环境

`MCST` 是用[Go](https://go.dev)编写的.

前提条件:

- [Task](https://taskfile.dev/installation)
- [Go 1.22+](https://go.dev/doc/install)
- [Aria2](https://github.com/aria2/aria2/releases)

在任意位置克隆`MCST`:

```shell
git clone git@github.com:Arama0517/MCST.git
```

使用`cd`进入目录并下载依赖:

```shell
task setup
```

确保一切正常的好办法是运行测试:

```shell
task test
```

## 测试你的更改

你可以为你的更改创建一个分支, 并尝试从源代码开始构建:

```shell
task build
```

当你对更改满意时, 我们建议你运行:

> 此操作需要安装[golangci-lint](https://golangci-lint.run/welcome/install)

```shell
task ci
```

在提交更改之前, 我们还建议你运行:

> 此操作需要安装[gofumpt](https://github.com/mvdan/gofumpt)

```shell
task fmt
```

## 创建提交

提交信息应该格式良好, 为了使其"标准化", 我们使用 `Conventional Commits` 格式

你可以关注[他们网站上的文档](https://www.conventionalcommits.org/)

## 提交拉取请求(Pull request)

将你的分支推送至你的`MCST`fork(分叉)并对main分支打开拉取请求
