package tcp_ip_proxys

// import (
// 	"io"
// 	"log"
// 	"net"
// 	"sync"
// 	"syscall"

// 	"github.com/alex60217101990/proxy.git/external/consts"
// 	"github.com/alex60217101990/proxy.git/external/enums"
// 	"github.com/alex60217101990/proxy.git/external/logger"
// 	"github.com/alex60217101990/proxy.git/external/models"
// 	"go.uber.org/zap"
// )

// type ProxyCore struct {
// 	laddr, raddr interface{}
// 	lconn, rconn io.ReadWriteCloser
// 	configs      *models.SingleProxyAddrConfigs
// 	closeCh      chan struct{}
// 	exitCh       chan interface{}
// 	clientConns  sync.Map
// 	bufferPool   sync.Pool
// }

// func (p *ProxyCore) Init() {
// 	syscall.Unlink(consts.ProxySocket)
// 	log.Println("Starting echo server")
// 	ln, err := net.Listen("unix", consts.ProxySocket)
// }

// // NewProxyCore - Create a new Proxy instance with or without a remote TLS server for
// // which we want to unwrap the TLS to be able to connect without encryption
// // locally. Takes over local connection passed in,
// // and closes it when finished.
// func NewProxyCore(configs *models.SingleProxyAddrConfigs) Proxy {
// 	if configs == nil {
// 		logger.Sugar.Fatal("proxy server configs is empty")
// 	}
// 	proxy := &ProxyCore{
// 		configs:    configs,
// 		closeCh:    make(chan struct{}, 1),
// 		exitCh:     make(chan interface{}, 1),
// 		bufferPool: sync.Pool{New: func() interface{} { return make([]byte, configs.Conn.MaxBufferSize) }},
// 	}
// 	var err error
// 	switch configs.Protocol {
// 	case enums.TCP:
// 		proxy.laddr, err = net.ResolveTCPAddr("tcp", configs.LisAddress)
// 		proxy.raddr, err = net.ResolveTCPAddr("tcp", configs.RecAddress)
// 	// 	// set lconn configs
// 	// if proxy.configs.Conn != nil {
// 	// 	if conn, ok := p.lconn.(TCPConnection); ok {
// 	// 		conn.SetNoDelay(true)
// 	// 	}
// 	// }
// 	case enums.UDP:
// 		proxy.laddr, err = net.ResolveUDPAddr("udp", configs.LisAddress)
// 		proxy.raddr, err = net.ResolveUDPAddr("udp", configs.RecAddress)
// 	}
// 	if err != nil {
// 		logger.Sugar.Fatalf("parse UDP/TCP listener address: [%v]\n", err)
// 	}

// 	return proxy
// }

// func (p *ProxyCore) Listen() {
// 	var (
// 		err error
// 	)
// 	defer func() {
// 		if r := recover(); r != nil {
// 			logger.Sugar.Fatal(err)
// 		}
// 		p.lconn.Close()
// 	}()
// 	switch p.configs.Protocol {
// 	case enums.TCP:
// 		//conn, err = net.Core("tcp", laddr)
// 	case enums.UDP:
// 		p.lconn, err = net.ListenUDP("udp", p.laddr.(*net.UDPAddr))
// 		if err != nil {
// 			logger.Sugar.Fatal(err)
// 		}
// 		go func() {
// 			for {
// 				//...

// 			}
// 		}()
// 	}

// }

// func (p *ProxyCore) UDPServerPipe() <-chan []byte {
// 	msgCh := make(chan []byte)
// 	go func() {
// 		defer close(msgCh)
// 		// defer p.bufferPool.Put(msg)
// 		msg := p.bufferPool.Get().([]byte)
// 		size, srcAddress, err := p.lconn.ReadFromUDP(msg[0:])
// 		if err != nil {
// 			p.Logger.Error(zap.Error(err))
// 			return
// 		}
// 		if p.configs.IsDebug {
// 			logger.Magenta.Printf("listener packet-received: bytes=%d from=%s\n", n, addr.String())
// 		}
// 		msgCh <- msg
// 	}()
// 	return msgCh
// }

// func (p *ProxyCore) UDPClientPipe() <-chan []byte {
// 	msgCh := make(chan []byte)
// 	go func() {
// 		defer close(msgCh)
// 		msg := p.bufferPool.Get().([]byte)
// 		size, srcAddress, err := p.rconn.ReadFromUDP(msg[0:])
// 		if err != nil {
// 			p.Logger.Error(zap.Error(err))
// 			return
// 		}
// 		if p.configs.IsDebug {
// 			logger.Cyan.Printf("upstreaam packet-received: bytes=%d from=%s\n", n, addr.String())
// 		}
// 		msgCh <- msg
// 	}()
// 	return msgCh
// }
