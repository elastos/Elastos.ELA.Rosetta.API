Elastos ELA Rosetta API
===========
[![GoDoc](https://godoc.org/github.com/elastos/Elastos.ELA.Coinbase.API?status.svg)](https://github.com/coinbase/rosetta-sdk-go) [![Build Status](https://travis-ci.com/elastos/Elastos.ELA.Coinbase.API.svg?branch=release_v0.0.1)](https://travis-ci.com/elastos/Elastos.ELA.Coinbase.API) [![Go Report Card](https://goreportcard.com/badge/github.com/elastos/Elastos.ELA.Coinbase.API)](https://goreportcard.com/report/github.com/elastos/Elastos.ELA.Coinbase.API) [![License](https://badges.fyi/github/license/elastos/Elastos.ELA.Coinbase.API)](https://github.com/elastos/Elastos.ELA.Coinbase.API/blob/master/LICENSE) [![Latest tag](https://badges.fyi/github/latest-tag/elastos/Elastos.ELA.Coinbase.API)](https://github.com/elastos/Elastos.ELA.Coinbase.API/releases)
## Introduction
A complete Rosetta Server implementation for [Elastos](https://github.com/elastos/Elastos.ELA).

## Build and run step by step

### 1. macOS Prerequisites

Make sure the [macOS version](https://en.wikipedia.org/wiki/MacOS#Release_history) Mojave or later 64-bit Intel.

```bash
$ uname -srm
Darwin 18.7.0 x86_64
```

Use [Homebrew](https://brew.sh/) to install Golang 1.16

```bash
$ brew install go@1.16
```

Check the golang version. Make sure they are the following version number or above.

```bash
$ go version
go version go1.16.5 darwin/amd64
```

### 2. Ubuntu Prerequisites

Make sure your ubuntu version is 18.04 or later.

```bash
$ cat /etc/issue
Ubuntu 18.04.5 LTS \n \l
```

Install Git.

```bash
$ sudo apt-get install -y git
```

Install Go distribution.

```bash
$ curl -O https://golang.org/dl/go1.13.15.linux-amd64.tar.gz
$ tar -xvf go1.13.15.linux-amd64.tar.gz
$ sudo chown -R root:root ./go
$ sudo mv go /usr/local
$ export GOPATH=$HOME/go
$ export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
$ source ~/.profile
```

### 3. Clone source code
Make sure you are in the working folder.
```bash
$ git clone https://github.com/elastos/Elastos.ELA.Rosetta.API.git
```

If clone works successfully, you should see folder structure like Elastos.ELA.Rosetta.API/Makefile

### 4. Run directly
1. Run `make runserver`
2. Run `make runclient` (in a new terminal window)
3. Run `make runfetcher` (in a new terminal window)

### 5. Build and run
1.build all
 ```shell script
make all
 ```
2.Run server
```shell script
./rs-server
```
3.Run client
```shell script
./rs-client
```
4.Run fetcher
```shell script
./rs-fetcher
```