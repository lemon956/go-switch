markdown
# go-switch

![Release](https://img.shields.io/github/v/release/lemon9563306/go-switch)

<p align="center">
  <a href="README_EN.md">英文</a>
</p>

## 目录

- [go-switch](#go-switch)
  - [目录](#目录)
  - [项目背景](#项目背景)
  - [功能](#功能)
  - [特性](#特性)
    - [\* 跨平台](#-跨平台)
    - [\* 安装方便](#-安装方便)
  - [安装](#安装)
    - [下载压缩包](#下载压缩包)
  - [最后](#最后)

## 项目背景
    go-switch 是一个可以跨平台进行golang版本管理工具,

## 功能
    help    - 显示此帮助信息
    install - 安装指定的 Go 版本
    switch  - 选择 Go 版本
    list    - 列出所有已安装的 Go 版本
    listall - 列出所有可用的 Go 版本
    delete  - 删除指定的 Go 版本
    lint    - 管理 golangci-lint v1/v2

### golangci-lint 管理
    goswitch lint install v1         - 安装最新的 golangci-lint v1
    goswitch lint install v2         - 安装最新的 golangci-lint v2
    goswitch lint install v1 v1.64.8 - 安装指定的 golangci-lint v1 版本
    goswitch lint switch v1          - 切换当前使用的 golangci-lint 到 v1
    goswitch lint switch v2          - 切换当前使用的 golangci-lint 到 v2
    goswitch lint list               - 列出已安装的 golangci-lint 版本
    goswitch lint env                - 查看 golangci-lint 环境信息

    golangci-lint 的当前入口复用 GOPATH/bin，不需要单独添加 lint 环境变量。
    不同版本的二进制文件保存在 go-switch 的 tools/golangci-lint 目录中。


## 特性
### * 跨平台
    - windows
    - linux
    - macos
### * 安装方便

## 安装
### 下载压缩包
    linux | macos:
        解压缩后将二进制文件移动到 /usr/local/bin 目录下
        expample cmd: sudo mv ./go-switch /usr/local/bin
    windows:
        创建一个目录,将下载的压缩文件解压缩到创建的目录,并将目录添加到环境变量中
    
## 最后
    如果有需求请提出issues 
