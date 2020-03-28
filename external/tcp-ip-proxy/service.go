package tcp_ip_proxy

import (
	"log"
	"net"
)

type ProxyCore struct {
}

func (p *ProxyCore) Init() {
	syscall.unlink(consts.ProxySocket)
	log.Println("Starting echo server")
	ln, err := net.Listen("unix")
}
