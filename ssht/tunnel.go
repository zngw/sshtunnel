// @Title SSH隧道
// @Description $
// @Author  55
// @Date  2021/11/7
package ssht

import (
	"fmt"
	"github.com/zngw/log"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// SSH隧道数据
type sshTunnel struct {
	host    string            // ssh服务器 地址:端口 如： 127.0.0.1:22
	user    string            // ssh登录用户名
	pwd     string            // ssh登录密码
	pkey    string            // ssh密钥
	client  *ssh.Client       // ssh链接
	tunnels map[string]string // 转发 key: local->remote, value: 监听`地址:端口`
}

// 链接ssh
func (t *sshTunnel) connect() (err error) {
	if t.client != nil {
		// 已经连接了，不重新连接
		return
	}

	// 验证数据
	var auth []ssh.AuthMethod
	auth = make([]ssh.AuthMethod, 0)
	if t.pkey != "" {
		var key ssh.Signer = nil
		if len(t.pkey) >= 1024 {
			// 如果密钥大于1024字符，优先考虑是否为密钥字符串
			key, err = ssh.ParsePrivateKey([]byte(t.pkey))
			if err != nil {
				key = nil
			}
		}

		if key == nil {
			// 读取密钥文件
			pKeyBytes, err1 := ioutil.ReadFile(t.pkey)
			if err1 != nil {
				err = fmt.Errorf("读取[%s]密钥文件%s错误，%v", t.host, t.pkey, err1)
				return
			}

			key, err1 = ssh.ParsePrivateKey(pKeyBytes)
			if err1 != nil {
				err = fmt.Errorf("解析[%s]密钥文件%s错误，%v", t.host, t.pkey, err1)
				return
			}
		}

		auth = append(auth, ssh.PublicKeys(key))
	} else if t.pwd != "" {
		auth = append(auth, ssh.Password(t.pwd))
	}

	//创建sshp登陆配置
	cfg := &ssh.ClientConfig{
		Timeout:         30 * time.Second,            //连接超时时间
		User:            t.user,                      // 账号
		Auth:            auth,                        // 验证
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //接受任何主机密钥
	}

	t.client, err = ssh.Dial("tcp", t.host, cfg)
	if err != nil {
		err = fmt.Errorf("连接[%s]服务器错误，%v", t.host, err)
		return
	}

	log.Info("tunnel", "连接[%s]服务器成功", t.host)

	return
}

// 断开连接
func (t *sshTunnel) disconnect() {
	if t.client != nil {
		_ = t.client.Close()
	}
}

// 隧道转发
// remote 开启隧道的远程主机配置，格式为【IP地址:端口】
// local 开启隧道映射到本地的配置，格式为【IP地址:端口】, 若端口为0或空，自动分配50000-59999间的端口
func (t *sshTunnel) forward(remote, local string) (listen string, err error) {
	// 隧道Id
	tid := fmt.Sprintf("%s->%s", local, remote)

	if t, ok := t.tunnels[tid]; ok {
		// 如果已经存在对应监听，直接返回
		listen = t
		return
	}

	// 本地监听 local listener
	var ll net.Listener

	s := strings.Split(local, ":")
	if len(s) != 2 {
		err = fmt.Errorf("[%s] 本地监听地址错误，%s", tid, local)
		return
	}

	// 随机端口[50000,59999]
	port, err := strconv.Atoi(s[1])
	if err != nil || port == 0 {
		port = rand.Intn(10000) + 50000
	}

	// 尝试55次
	for i := 0; i < 55; i++ {
		listen = fmt.Sprintf("%s:%d", s[0], port)
		ll, err = net.Listen("tcp", listen)
		if err != nil {
			// 随机端口重试
			port = rand.Intn(10000) + 50000
		} else {
			break
		}
	}

	if ll == nil {
		// 超出尝试连接次数
		err = fmt.Errorf("[%s] 本地监听开始失败，%v", tid, err)
		listen = ""
		return
	}

	log.Info("tunnel", "[%s] 监听开启: %s", tid, listen)
	go t.accept(tid, ll, remote)

	return
}

// 启用监听
func (t *sshTunnel) accept(tid string, ll net.Listener, remote string) {
	defer func() {
		_ = ll.Close()
		log.Info("tunnel", "[%s] 监听关闭", tid)
	}()

	sid := int64(1)
	for {
		// 本地监听接收 local Conn
		lc, err := ll.Accept()
		if err != nil {
			log.Error("tunnel", "[%s] 接收连接失败, %v", tid, err)
			return
		}

		// 连接ssh
		err = t.connect()
		if err != nil {
			log.Error("tunnel", "[%s] %v", tid, err)
			_ = lc.Close()
			continue
		}

		// 获取远程连接 remote Conn
		rc, err := t.client.Dial("tcp", remote)
		if err != nil {
			log.Error("tunnel", "[%s] 获取远程连接失败, %v", tid, err)
			_ = t.client.Close()
			_ = lc.Close()
			continue
		}

		if sid >= math.MaxInt64 {
			sid = 0
		}
		sid++
		cid := fmt.Sprintf("%s:%d", tid, sid)
		go t.transfer(cid, lc, rc)
	}
}

// 消息转发
func (t *sshTunnel) transfer(cid string, lc, rc net.Conn) {
	defer rc.Close()
	defer lc.Close()
	go func() {
		defer lc.Close()
		defer rc.Close()
		_, _ = io.Copy(rc, lc)
	}()
	log.Info("tunnel", "[%s] 通道已连接!", cid)
	_, _ = io.Copy(lc, rc)
	log.Info("tunnel", "[%s] 通道已断开!", cid)
}
