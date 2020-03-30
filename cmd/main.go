package main

import (
	"github.com/alex60217101990/proxy.git/external/tcp-ip-proxys"
)

func main() {
	p := tcp_ip_proxys.ProxyCore{}
	p.Listen()
}
