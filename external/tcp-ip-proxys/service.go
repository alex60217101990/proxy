package tcp_ip_proxys

import (
	"io"
	"log"
	"net"
	"syscall"

	"github.com/alex60217101990/proxy.git/external/consts"
	"github.com/alex60217101990/proxy.git/external/enums"
	"github.com/alex60217101990/proxy.git/external/logger"
	"github.com/alex60217101990/proxy.git/external/models"
)

type ProxyCore struct {
	laddr, raddr interface{}
	lconn, rconn io.ReadWriteCloser
	configs      *models.SingleProxyAddrConfigs
	closeCh, exitCh chan struct{}
	bufferPool             sync.Pool
}

func (p *ProxyCore)p *ProxyCore Init() {
	syscall.Unlink(consts.ProxySocket)
	log.Println("Starting echo server")
	ln, err := net.Listen("unix", consts.ProxySocket)
}

// NewProxyCore - Create a new Proxy instance with or without a remote TLS server for
// which we want to unwrap the TLS to be able to connect without encryption
// locally. Takes over local connection passed in,
// and closes it when finished.
func NewProxyCore(configs *models.SingleProxyAddrConfigs) Proxy {
	if configs == nil {
		logger.Sugar.Fatal("proxy server configs is empty")
	}
	proxy := &ProxyCore{
		configs: configs,
		closeCh: make(chan struct{}, 1),
		exitCh: make(chan struct{}, 1),
	}
	var err error
	switch configs.Protocol {
	case enums.TCP:
		proxy.laddr, err = net.ResolveTCPAddr("tcp", configs.LisAddress)
		proxy.raddr, err = net.ResolveTCPAddr("tcp", configs.RecAddress)
	// 	// set lconn configs
	// if proxy.configs.Conn != nil {
	// 	if conn, ok := p.lconn.(TCPConnection); ok {
	// 		conn.SetNoDelay(true)
	// 	}
	// }
	case enums.UDP:
		proxy.laddr, err = net.ResolveUDPAddr("udp", configs.LisAddress)
		proxy.raddr, err = net.ResolveUDPAddr("udp", configs.RecAddress)
	}
	if err = nil {
		logger.Sugar.Fatalf("parse UDP/TCP listener address: [%v]\n", err)
	}
	
	return proxy
}

func (p *ProxyCore) Listen() {
	var (
		err error
	)
	defer func ()  {
		if r := recover(); r != nil {
			logger.Sugar.Debug("packet received",
			zap.String("src address", packetSourceString),
			zap.Int("src port", pa.src.Port),
			zap.String("packet", string(pa.data)),
			zap.Int("size", len(pa.data)),
		)
		}
	}
	switch configs.Protocol {
	case enums.TCP:
		conn, err = net.Core("tcp", laddr)
	case enums.UDP:
		conn, err = net.Listen("udp", p.laddr)
	}
	
}
