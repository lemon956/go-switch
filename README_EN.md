# go-switch

![Release](https://img.shields.io/github/v/release/lemon9563306/go-switch)

<p align="center">
    <a href="README.md">Chinese</a>
</p>

## Table of Contents

- [go-switch](#go-switch)
  - [Table of Contents](#table-of-contents)
  - [Project Background](#project-background)
  - [Features](#features)
  - [Characteristics](#characteristics)
    - [\* Cross-Platform](#-cross-platform)
    - [\* Easy Installation](#-easy-installation)
  - [Installation](#installation)
    - [Download the Archive](#download-the-archive)
  - [Conclusion](#conclusion)

## Project Background
    go-switch is a cross-platform tool for managing Golang versions.

## Features
    help    - Display this help information
    install - Install the specified Go version
    switch  - Select the Go version
    list    - List all installed Go versions
    listall - List all available Go versions
    delete  - Delete the specified Go version
    lint    - Manage golangci-lint v1/v2

### golangci-lint management
    goswitch lint install v1         - Install the latest golangci-lint v1
    goswitch lint install v2         - Install the latest golangci-lint v2
    goswitch lint install v1 v1.64.8 - Install a specific golangci-lint v1 version
    goswitch lint switch v1          - Switch current golangci-lint to v1
    goswitch lint switch v2          - Switch current golangci-lint to v2
    goswitch lint list               - List installed golangci-lint versions
    goswitch lint env                - Show golangci-lint environment

    The current golangci-lint entry reuses GOPATH/bin, so no separate lint PATH is required.
    Versioned binaries are stored under go-switch's tools/golangci-lint directory.

## Characteristics
### * Cross-Platform
    - windows
    - linux
    - macos
### * Easy Installation

## Installation
### Download the Archive
    linux | macos:
        Extract the archive and move the binary file to the /usr/local/bin directory
        Example cmd: sudo mv ./go-switch /usr/local/bin
    windows:
        Create a directory, extract the downloaded archive to the created directory, and add the directory to the environment variable

## Conclusion
    If you have any requests, please raise an issue.
