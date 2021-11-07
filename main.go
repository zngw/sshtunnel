// @Title
// @Description $
// @Author  55
// @Date  2021/11/05
package main

import (
	"flag"
	"fmt"
	"github.com/zngw/log"
	"os"
	"os/signal"
	"sshtunnel/config"
	"sshtunnel/ssht"
	"strings"
	"syscall"
)

func main() {
	// 读取命令行配置文件参数
	c := flag.String("c", "./config.yaml", "默认配置为 ./config.yaml")
	l := flag.String("l", "", "格式如：remote:port1:local:port2")
	h := flag.String("h", "", "格式如：ssh://root:123456@192.168.1.55:22")
	p := flag.String("p", "", "密钥文件路径，密码登录此参数可忽略")
	flag.Parse()

	if len(*l) > 0 {
		s := strings.Split(*l, ":")
		if len(s) != 4 {
			log.Info("net", "传数参数错误")
			return
		}

		listen, err := ssht.TunnelUriByKey(*h, *p, fmt.Sprintf("%s:%s", s[0], s[1]), fmt.Sprintf("%s:%s", s[2], s[3]))
		if err != nil {
			log.Info("net", "启用隧道失败 %v", err)
		} else {
			log.Info("net", "启用隧道成功 %s", listen)
		}
	} else if len(*c) > 0 {
		// 存在配置文件
		err := config.Init("./config.yaml")
		if err != nil {
			log.Error("main", "%v", err)
			return
		}

		for _, s := range config.Config {
			for _, t := range s.Tunnels {
				listen, err := ssht.TunnelUriByKey(s.Uri, s.Pkey, t.Remote, t.Local)
				if err != nil {
					log.Info("main", "启用隧道失败 %v", err)
				} else {
					log.Info("main", "启用隧道成功 %s", listen)
				}
			}
		}
	}

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-signalChan
	ssht.CloseAllTunnel()
	os.Exit(0)
}
