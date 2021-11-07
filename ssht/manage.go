// @Title SSH隧道管理
// @Description $
// @Author  55
// @Date  2021/11/7
package ssht

import (
	"fmt"
	"math/rand"
	"net/url"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// SSH服务器数据，key-ssh host
var tunnelMap = make(map[string]*sshTunnel)

// 以URI格式添加隧道
// uri 格式: ssh://root:123456@192.168.1.55:22
// remote 开启隧道的远程主机配置，格式为【IP地址:端口】
// local 开启隧道映射到本地的配置，格式为【IP地址:端口】, 若端口为0或空，自动分配50000-59999间的端口
// 返回本地监听`地址:端口`和错误, 如果 `local != listen`说明自动分配过端口
func TunnelUri(uri, remote, local string) (listen string, err error) {
	return TunnelUriByKey(uri, "", remote, local)
}

// 以URI格式加密钥认证方式添加隧道
// uri 格式: ssh://root:123456@192.168.1.55:22
// pkey 密钥文件路径 或 密钥字符串
// remote 开启隧道的远程主机配置，格式为【IP地址:端口】
// local 开启隧道映射到本地的配置，格式为【IP地址:端口】, 若端口为0或空，自动分配50000-59999间的端口
// 返回本地监听`地址:端口`和错误, 如果 `local != listen`说明自动分配过端口
func TunnelUriByKey(uri, pkey, remote, local string) (listen string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		err = fmt.Errorf("解析 %s 失败，%v", uri, err)
		return
	}

	pwd, _ := u.User.Password()
	return Tunnel(u.Host, u.User.Username(), pwd, pkey, remote, local)
}

// 账号密码格式添加隧道
// host ssh服务器 `地址:端口`
// user 登录用户名
// pwd 登录密码
// remote 开启隧道的远程主机配置，格式为【IP地址:端口】
// local 开启隧道映射到本地的配置，格式为【IP地址:端口】, 若端口为0或空，自动分配50000-59999间的端口
// 返回本地监听`地址:端口`和错误, 如果 `local != listen`说明自动分配过端口
func TunnelByPassword(host, user, pwd, remote, local string) (listen string, err error) {
	return Tunnel(host, user, pwd, "", remote, local)
}

// 密钥格式添加隧道
// host ssh服务器 `地址:端口`
// user 登录用户名
// pkey 密钥文件
// remote 开启隧道的远程主机配置，格式为【IP地址:端口】
// local 开启隧道映射到本地的配置，格式为【IP地址:端口】, 若端口为0或空，自动分配50000-59999间的端口
// 返回本地监听`地址:端口`和错误, 如果 `local != listen`说明自动分配过端口
func TunnelByKey(host, user, pkey, remote, local string) (listen string, err error) {
	return Tunnel(host, user, "", pkey, remote, local)
}

// 添加隧道转发
// host ssh服务器 `地址:端口`
// user 登录用户名
// pwd 登录密码
// pkey 密钥文件
// remote 开启隧道的远程主机配置，格式为【IP地址:端口】
// local 开启隧道映射到本地的配置，格式为【IP地址:端口】, 若端口为0或空，自动分配50000-59999间的端口
// 返回本地监听`地址:端口`和错误, 如果 `local != listen`说明自动分配过端口
func Tunnel(host, user, pwd, pkey, remote, local string) (listen string, err error) {
	ts := getTunnel(host)
	if ts == nil {
		ts = &sshTunnel{
			host:   host,
			user:   user,
			pwd:    pwd,
			pkey:   pkey,
			client: nil,
		}

		err = ts.connect()
		if err != nil {
			return
		}
		tunnelMap[host] = ts
	}

	listen, err = ts.forward(remote, local)
	return
}

// 关闭所有隧道
func CloseAllTunnel() {
	for _, v := range tunnelMap {
		v.disconnect()
	}
}

// 获取SSH服务器数据
func getTunnel(host string) (ts *sshTunnel) {
	if t, ok := tunnelMap[host]; ok {
		ts = t
		return
	}

	return
}
