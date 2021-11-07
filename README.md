[comment]: <> (dtapps)
[![GitHub Org's stars](https://img.shields.io/github/stars/zngw)](https://github.com/zngw)

[comment]: <> (go)
[![godoc](https://pkg.go.dev/badge/github.com/zngw/sshtunnel?status.svg)](https://pkg.go.dev/github.com/zngw/sshtunnel)
[![oproxy.cn](https://goproxy.cn/stats/github.com/zngw/sshtunnel/badges/download-count.svg)](https://goproxy.cn/stats/github.com/zngw/sshtunnel)
[![goreportcard.com](https://goreportcard.com/badge/github.com/zngw/sshtunnel)](https://goreportcard.com/report/github.com/zngw/sshtunnel)
[![deps.dev](https://img.shields.io/badge/deps-go-red.svg)](https://deps.dev/go/github.com%2Fdtapps%2Fgo-ssh-tunnel)

[comment]: <> (github.com)
[![watchers](https://badgen.net/github/watchers/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/watchers)
[![stars](https://badgen.net/github/stars/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/stargazers)
[![forks](https://badgen.net/github/forks/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/network/members)
[![issues](https://badgen.net/github/issues/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/issues)
[![branches](https://badgen.net/github/branches/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/branches)
[![releases](https://badgen.net/github/releases/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/releases)
[![tags](https://badgen.net/github/tags/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/tags)
[![license](https://badgen.net/github/license/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/blob/master/LICENSE)
[![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/zngw/sshtunnel)](https://github.com/zngw/sshtunnel)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/releases)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/tags)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/pulls)
[![GitHub issues](https://img.shields.io/github/issues/zngw/sshtunnel)](https://github.com/zngw/sshtunnel/issues)
[![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/zngw/sshtunnel)](https://github.com/zngw/sshtunnel)
[![GitHub language count](https://img.shields.io/github/languages/count/zngw/sshtunnel)](https://github.com/zngw/sshtunnel)
[![GitHub search hit counter](https://img.shields.io/github/search/zngw/sshtunnel/go)](https://github.com/zngw/sshtunnel)
[![GitHub top language](https://img.shields.io/github/languages/top/zngw/sshtunnel)](https://github.com/zngw/sshtunnel)
# sshtunnel

go语言实现的一个ssh端口映射的程序

# 使用二进制程序

可以自己下程序自行编译或去 https://github.com/zngw/sshtunnel/releases 下载已经编译好的对应版本

## 单端口映射使用

```bash
sshtunnel -l 127.0.0.1:3717:0.0.0.0:3717 -h ssh://root:123456@192.168.1.55:22
```

### 说明

* -l: 映射端口。127.0.0.1：3717 开启隧道的远程主机配置；0.0.0.0:3717 开启隧道映射到本地的配置
* -h: ssh连接主机参数，uri格式
* -p: 密钥文件，若是密码登录此项可忽略

## 使用配置文件

```bash
sshtunnel -c ./config.yaml
```

### 配置说明

配置文件为`yaml`格式,例如：

```yaml
-
  # SSH连接配置，uri格式
  uri: ssh://root@192.168.1.55:22
  # 密钥文件，若是密码登录此项可忽略
  pkey: ./rsa.pem
  # 隧道转发配置
  tunnels:
    -
      # 开启隧道的远程主机配置，格式为【IP地址:端口】
      remote: 127.0.0.1:27017
      # 开启隧道映射到本地的配置，格式为【IP地址:端口】
      local: 0.0.0.0:27017
    -
      remote: 127.0.0.1:6379
      local: 0.0.0.0:6379
-
  # SSH连接配置，uri格式
  uri: ssh://root:123456@192.168.1.56:22
  # 隧道转发配置
  tunnels:
    -
      # 开启隧道的远程主机配置，格式为【IP地址:端口】
      remote: 127.0.0.1:27017
      # 开启隧道映射到本地的配置，格式为【IP地址:端口】
      local: 0.0.0.0:27017
```