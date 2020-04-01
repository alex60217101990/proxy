package tcp_ip_proxys

// import (
// 	"fmt"
// 	"net"
// 	"runtime"
// 	"sync"
// 	"time"

// 	"github.com/alex60217101990/proxy.git/external/logger"
// 	"github.com/alex60217101990/proxy.git/external/models"
// 	"go.uber.org/zap"
// )

// type UDPProxy struct {
// 	BindPort               int
// 	BindAddress            string
// 	UpstreamAddress        string
// 	UpstreamPort           int
// 	Debug                  bool
// 	listenerConn           *net.UDPConn
// 	client                 *net.UDPAddr
// 	upstream               *net.UDPAddr
// 	BufferSize             int
// 	ConnTimeout            time.Duration
// 	ResolveTTL             time.Duration
// 	connsMap               sync.Map
// 	closed                 bool
// 	clientMessageChannel   chan models.UDPPacket
// 	upstreamMessageChannel chan models.UDPPacket
// 	bufferPool             sync.Pool
// }

// func NewUDPProxy(debug bool, bindPort int, bindAddress string, upstreamAddress string, upstreamPort int, bufferSize int, connTimeout time.Duration, resolveTTL time.Duration) Proxy {
// 	return &UDPProxy{
// 		Debug:                  debug,
// 		BindPort:               bindPort,
// 		BindAddress:            bindAddress,
// 		BufferSize:             bufferSize,
// 		ConnTimeout:            connTimeout,
// 		UpstreamAddress:        upstreamAddress,
// 		UpstreamPort:           upstreamPort,
// 		closed:                 false,
// 		ResolveTTL:             resolveTTL,
// 		clientMessageChannel:   make(chan models.UDPPacket, 50),
// 		upstreamMessageChannel: make(chan models.UDPPacket, 50),
// 		bufferPool:             sync.Pool{New: func() interface{} { return make([]byte, bufferSize) }},
// 	}
// }

// func (p *UDPProxy) clientConnectionReadLoop(clientAddr *net.UDPAddr, upstreamConn *net.UDPConn) {
// 	clientAddrString := clientAddr.String()
// 	for {
// 		msg := p.bufferPool.Get().([]byte)
// 		size, _, err := upstreamConn.ReadFromUDP(msg[0:])
// 		if err != nil {
// 			upstreamConn.Close()
// 			p.connsMap.Delete(clientAddrString)
// 			return
// 		}
// 		p.upstreamMessageChannel <- models.UDPPacket{
// 			Src:  clientAddr,
// 			Data: msg[:size],
// 		}
// 	}
// }

// func (p *UDPProxy) handlerUpstreamPackets() {
// 	for pa := range p.upstreamMessageChannel {
// 		logger.Sugar.Debug("forwarded data from upstream", zap.Int("size", len(pa.Data)), zap.String("data", string(pa.Data)))
// 		p.listenerConn.WriteTo(pa.Data, pa.Src)
// 		p.bufferPool.Put(pa.Data)
// 	}
// }

// func (p *UDPProxy) handleClientPackets() {
// 	for pa := range p.clientMessageChannel {
// 		packetSourceString := pa.Src.String()
// 		logger.Sugar.Debug("packet received",
// 			zap.String("src address", packetSourceString),
// 			zap.Int("src port", pa.Src.Port),
// 			zap.String("packet", string(pa.Data)),
// 			zap.Int("size", len(pa.Data)),
// 		)

// 		conn, found := p.connsMap.Load(packetSourceString)
// 		if !found {
// 			conn, err := net.ListenUDP("udp", p.client)
// 			logger.Sugar.Debug("new client connection",
// 				zap.String("local port", conn.LocalAddr().String()),
// 			)

// 			if err != nil {
// 				logger.Sugar.Error("upd proxy failed to dial", zap.Error(err))
// 				return
// 			}
// 			conn.SetDeadline(time.Now().Add(p.ConnTimeout))
// 			p.connsMap.Store(packetSourceString, conn)

// 			conn.WriteTo(pa.Data, p.upstream)
// 			go p.clientConnectionReadLoop(pa.Src, conn)
// 		} else {
// 			conn.(*net.UDPConn).WriteTo(pa.Data, p.upstream)
// 		}
// 		p.bufferPool.Put(pa.Data)
// 	}
// }

// func (p *UDPProxy) readLoop() {
// 	for !p.closed {
// 		msg := p.bufferPool.Get().([]byte)
// 		size, srcAddress, err := p.listenerConn.ReadFromUDP(msg[0:])
// 		if err != nil {
// 			logger.Sugar.Error("error", zap.Error(err))
// 			continue
// 		}
// 		p.clientMessageChannel <- models.UDPPacket{
// 			Src:  srcAddress,
// 			Data: msg[:size],
// 		}
// 	}
// }

// func (p *UDPProxy) resolveUpstreamLoop() {
// 	for !p.closed {
// 		if p.ResolveTTL > 0 {
// 			select {
// 			case <-time.After(p.ResolveTTL):
// 				upstreamAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.UpstreamAddress, p.UpstreamPort))
// 				if err != nil {
// 					logger.Sugar.Error("resolve error", zap.Error(err))
// 					continue
// 				}
// 				if p.upstream.String() != upstreamAddr.String() {
// 					p.upstream = upstreamAddr
// 					logger.Sugar.Info("upstream addr changed", zap.String("upstreamAddr", p.upstream.String()))
// 				}
// 			}
// 			// continue
// 		}
// 		// upstreamAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.UpstreamAddress, p.UpstreamPort))
// 		// if err != nil {
// 		// 	logger.Sugar.Error("resolve error", zap.Error(err))
// 		// 	continue
// 		// }
// 		// if p.upstream.String() != upstreamAddr.String() {
// 		// 	p.upstream = upstreamAddr
// 		// 	logger.Sugar.Info("upstream addr changed", zap.String("upstreamAddr", p.upstream.String()))
// 		// }
// 	}
// }

// // Close stops the proxy
// func (p *UDPProxy) Close() {
// 	logger.Sugar.Warn("Closing proxy")
// 	p.closed = true
// 	p.connsMap.Range(func(k, conn interface{}) bool {
// 		conn.(*net.UDPConn).Close()
// 		return true
// 	})
// 	if p.listenerConn != nil {
// 		p.listenerConn.Close()
// 	}
// }

// // Start starts the proxy
// func (p *UDPProxy) Start() {
// 	logger.Sugar.Info("Starting proxy")

// 	ProxyAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.BindAddress, p.BindPort))
// 	if err != nil {
// 		logger.Sugar.Error("error resolving bind address", zap.Error(err))
// 		return
// 	}
// 	p.upstream, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.UpstreamAddress, p.UpstreamPort))
// 	if err != nil {
// 		logger.Sugar.Error("error resolving upstream address", zap.Error(err))
// 	}
// 	p.client = &net.UDPAddr{
// 		IP:   ProxyAddr.IP,
// 		Port: p.BindPort, // 0
// 		Zone: ProxyAddr.Zone,
// 	}
// 	p.listenerConn, err = net.ListenUDP("udp", ProxyAddr)
// 	if err != nil {
// 		logger.Sugar.Error("error listening on bind port", zap.Error(err))
// 		return
// 	}
// 	logger.Sugar.Info("UDP Proxy started!")
// 	if p.ConnTimeout.Nanoseconds() == 0 {
// 		logger.Sugar.Warn("be warned that running without timeout to clients may be dangerous")
// 	}
// 	if p.ResolveTTL.Nanoseconds() > 0 {
// 		go p.resolveUpstreamLoop()
// 	} else {
// 		logger.Sugar.Warn("not refreshing upstream addr")
// 	}
// 	for i := 0; i < runtime.NumCPU(); i++ {
// 		go p.readLoop()
// 		go p.handleClientPackets()
// 		go p.handlerUpstreamPackets()
// 	}
// }
