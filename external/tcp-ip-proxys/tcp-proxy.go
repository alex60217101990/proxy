package tcp_ip_proxys

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"
	"sync"
	"time"

	"github.com/alex60217101990/proxy.git/external/consts"
	"github.com/alex60217101990/proxy.git/external/enums"
	"github.com/alex60217101990/proxy.git/external/logger"
	"github.com/alex60217101990/proxy.git/external/models"
	"github.com/tevino/abool"
	"go.uber.org/zap"
)

// TCPProxy ...
type TCPProxy struct {
	Debug           bool
	connMap         sync.Map // models.PipeConnTCP
	blackList       sync.Map
	bufferPool      sync.Pool
	listenerConn    io.ReadWriteCloser // *net.TCPConn
	listenerAddr    *net.TCPAddr
	closed          *abool.AtomicBool
	addConnCh       chan models.ConnSignal
	delConnCh       chan models.ConnSignal
	listenerConfigs *models.ConnConfigs
}

// NewTCPProxy ...
func NewTCPProxy(options ...func(*TCPProxy) error) Proxy {
	proxy := &TCPProxy{
		closed:    abool.New(),
		addConnCh: make(chan models.ConnSignal, 50),
		delConnCh: make(chan models.ConnSignal, 50),
	}
	proxy.closed.UnSet()
	if len(options) < 1 {
		logger.Sugar.Fatal(errToSmallOptionsList(enums.TCP))
	}
	for _, op := range options {
		err := op(proxy)
		if err != nil {
			logger.Sugar.Fatal(err)
		}
	}
	return proxy
}

// SetConfigs ...
func SetConfigs(conf *models.ConnConfigs) func(*TCPProxy) error {
	return func(service *TCPProxy) (err error) {
		if conf == nil {
			return errEmptyProxyListenerConfigs(enums.TCP)
		} else if len(conf.LisAddress) == 0 {
			return errEmptyProxyListenerAddress(enums.TCP)
		}
		service.listenerConfigs = conf
		service.listenerAddr, err = net.ResolveTCPAddr("tcp", conf.LisAddress)
		if err != nil {
			logger.Sugar.Errorf("Failed to resolve local address: %s", err)
			return errResolveLocalAddress(err)
		}
		if conf.MaxBufferSize == 0 {
			conf.MaxBufferSize = consts.MaxBufferSize
		}
		service.bufferPool = sync.Pool{New: func() interface{} { return make([]byte, conf.MaxBufferSize) }}
		return err
	}
}

// SetIsDebug ...
func SetIsDebug(debug bool) func(*TCPProxy) error {
	return func(service *TCPProxy) error {
		service.Debug = debug
		return nil
	}
}

func (p *TCPProxy) delClientsConnectionsLoop() {
	go func() {
		for !p.closed.IsSet() {
			select {
			case sig, ok := <-p.delConnCh:
				if ok {
					var err error
					connection, present := p.connMap.Load(sig.ConnConfigs.LisAddress)
					if present {
						err = connection.(io.ReadWriteCloser).Close()
						if err != nil {
							logger.Sugar.Errorf("close client connection: [%v] error: %v", sig.ConnConfigs.LisAddress, err)
							continue
						}
						p.connMap.Delete(sig.ConnConfigs.LisAddress)
						logger.Cyan.Printf("Connection with address: [%s], action: [%v] success.\n", sig.ConnConfigs.LisAddress, sig.OperationType)
					}
					err = p.listenerConn.Close()
					if err != nil {
						logger.Sugar.Errorf("close proxy server connection: [%v] error: %v", p.listenerAddr.String(), err)
						continue
					}
				}
			}
		}
	}()
}

func (p *TCPProxy) addClientsConnectionsLoop() {
	go func() {
		for !p.closed.IsSet() {
			select {
			case sig, ok := <-p.addConnCh:
				if ok {
					_, present := p.connMap.Load(sig.ConnConfigs.LisAddress)
					if present {
						logger.Cyan.Printf("Connection with address: [%s] already exists.\n", sig.ConnConfigs.LisAddress)
						continue
					}
					listenerAddr, err := net.ResolveTCPAddr("tcp", sig.ConnConfigs.LisAddress)
					if err != nil {
						logger.Sugar.Error(zap.Error(err))
						continue
					}
					var (
						conn       io.ReadWriteCloser
						tlsConfigs *tls.Config
					)
					tlsConfigs = nil
					if sig.ConnConfigs.Protocol == enums.TLS {
						if sig.ConnConfigs.Creeds != nil {
							if len(sig.ConnConfigs.Creeds.RootCAs) > 0 {
								roots := x509.NewCertPool()
								ok := roots.AppendCertsFromPEM([]byte(sig.ConnConfigs.Creeds.RootCAs))
								if !ok {
									logger.Sugar.Error("failed to parse root certificate")
									continue
								}
								tlsConfigs = &tls.Config{RootCAs: roots}
							} else if len(sig.ConnConfigs.Creeds.PublicKey) > 0 && len(sig.ConnConfigs.Creeds.PrivateKey) > 0 {
								var cert tls.Certificate
								cert, err = tls.LoadX509KeyPair(sig.ConnConfigs.Creeds.PublicKey, sig.ConnConfigs.Creeds.PrivateKey)
								if err != nil {
									logger.Sugar.Errorf("load TLS keys error: %s", err)
									continue
								}
								tlsConfigs = &tls.Config{
									Certificates:       []tls.Certificate{cert},
									InsecureSkipVerify: true,
								}
							}
						}
						conn, err = tls.Dial("tcp", listenerAddr.String(), tlsConfigs)
					} else {
						conn, err = net.DialTCP("tcp", nil, listenerAddr)
					}
					if err != nil {
						logger.Sugar.Error(zap.Error(err))
						continue
					}
					p.pipeData(conn, p.listenerConn, sig.ConnConfigs)
					p.connMap.Store(sig.ConnConfigs.LisAddress, conn)
					logger.Cyan.Printf("Connection with address: [%s], action: [%v] success.\n", sig.ConnConfigs.LisAddress, sig.OperationType)
				}
			}
		}
	}()
}

// EmmitAddSignal ...
func (p *TCPProxy) EmmitAddSignal(signal *models.ConnSignal) {
	select {
	case p.addConnCh <- *signal:
	default:
		logger.Sugar.Warn("Emmit ADD new TCP proxy server connection failed.")
	}
}

// EmmitDelSignal ...
func (p *TCPProxy) EmmitDelSignal(signal *models.ConnSignal) {
	select {
	case p.delConnCh <- *signal:
	default:
		logger.Sugar.Warn("Emmit DELETE new TCP proxy server connection failed.")
	}
}

func (p *TCPProxy) setTCPFlags(connI interface{}, config *models.ConnConfigs) {
	conn, ok := connI.(TCPConnection)
	if ok {
		if config.UseKeepAlive {
			conn.SetKeepAlive(true)
		}
		if config.UseKeepAlive && config.KeepAliveTimeout > 0 {
			conn.SetKeepAlivePeriod(time.Second * time.Duration(config.KeepAliveTimeout))
		}
		if config.UseNagles {
			conn.SetNoDelay(config.UseNagles)
		}
		if config.LingerSec > 0 {
			conn.SetLinger(config.LingerSec)
		}
	}
}

func (p *TCPProxy) setNetFlags(connI interface{}, config *models.ConnConfigs) {
	conn, ok := connI.(NetConnection)
	if ok {
		if config.Deadline > 0 {
			conn.SetDeadline(time.Now().Add(time.Second * time.Duration(config.Deadline)))
		} else {
			conn.SetDeadline(time.Now().Add(time.Second * time.Duration(consts.ConnDeadline)))
		}
	}
}

func (p *TCPProxy) setNetReadFlag(connI interface{}, config *models.ConnConfigs) {
	conn, ok := connI.(NetConnection)
	if ok {
		if config.Deadline > 0 {
			conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(config.Deadline)))
		} else {
			conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(consts.ConnDeadline)))
		}
	}
}

func (p *TCPProxy) setNetWriteFlag(connI interface{}, config *models.ConnConfigs) {
	conn, ok := connI.(NetConnection)
	if ok {
		if config.Deadline > 0 {
			conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(config.Deadline)))
		} else {
			conn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(consts.ConnDeadline)))
		}
	}
}

func (p *TCPProxy) pipeIn(connIn, connOut io.ReadWriteCloser, config *models.ConnConfigs) {
	go func() {
		var (
			writeLen int64
			err      error
		)
		buffIn := p.bufferPool.Get().([]byte)
		defer func() {
			p.EmmitDelSignal(&models.ConnSignal{
				OperationType: enums.DELETE,
				ConnConfigs:   config,
			})
			p.bufferPool.Put(buffIn)
			logger.Cyan.Printf("close connection: [%v] success.\n", config.LisAddress)
		}()
		p.setNetFlags(connIn, config)
		p.setNetFlags(connOut, config)
		writeLen, err = io.CopyBuffer(connIn, connOut, buffIn)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		logger.Magenta.Printf("copy bytes from IN to OUT: [%d],\n", writeLen)
	}()
}

func (p *TCPProxy) pipeOut(connIn, connOut io.ReadWriteCloser, config *models.ConnConfigs) {
	go func() {
		var (
			writeLen int64
			err      error
		)
		buffOut := p.bufferPool.Get().([]byte)
		defer func() {
			p.EmmitDelSignal(&models.ConnSignal{
				OperationType: enums.DELETE,
				ConnConfigs:   config,
			})
			p.bufferPool.Put(buffOut)
			logger.Cyan.Printf("close connection: [%v] success.\n", config.LisAddress)
		}()
		p.setNetFlags(connIn, config)
		p.setNetFlags(connOut, config)
		writeLen, err = io.CopyBuffer(connOut, connIn, buffOut)
		if err != nil {
			logger.Sugar.Error(err)
			return
		}
		logger.Magenta.Printf("copy bytes from OUT to IN: [%d],\n", writeLen)
	}()
}

func (p *TCPProxy) pipeData(connL, connU io.ReadWriteCloser, config *models.ConnConfigs) {
	p.setTCPFlags(connL, config)
	// connections pipe
	if !p.closed.IsSet() {
		p.pipeIn(connL, connU, config)
		p.pipeOut(connL, connU, config)
	}
}

// Close stops the proxy
func (p *TCPProxy) Close() {
	var err error
	logger.Green.Printf("Closing proxy, type: [%s], addr: [%s].\n", p.listenerConfigs.Protocol, p.listenerAddr.String())
	p.closed.SetToIf(false, true)
	p.connMap.Range(func(k, conn interface{}) bool {
		if err = conn.(io.ReadWriteCloser).Close(); err == nil {
			logger.Cyan.Printf("closing connection to remote: [%v] success.\n", k)
		} else {
			logger.Red.Printf("closing connection to remote: [%v] with error: %v\n", k, err)
		}
		return true
	})
	if p.listenerConn != nil {
		p.listenerConn.Close()
	}
}

func (p *TCPProxy) detectBlackAddr() <-chan interface{} {
	resp := make(chan interface{}, 1)
	go func(){
		defer func(){
			close(resp)
		}()
		
	}()
	return resp
}

func (p *TCPProxy) Listen(addrList ...*models.ConnConfigs) {
	var tlsConfigs *tls.Config
	if p.listenerConfigs.Protocol == enums.TLS {
		if p.listenerConfigs.Creeds != nil {
			if len(p.listenerConfigs.Creeds.RootCAs) > 0 {
				roots := x509.NewCertPool()
				ok := roots.AppendCertsFromPEM([]byte(p.listenerConfigs.Creeds.RootCAs))
				if !ok {
					logger.Sugar.Fatal(errParseRootCert(p.listenerConfigs.Protocol, p.listenerConfigs.LisAddress))
				}
				tlsConfigs = &tls.Config{RootCAs: roots}
			} else if len(p.listenerConfigs.Creeds.PublicKey) > 0 && len(p.listenerConfigs.Creeds.PrivateKey) > 0 {
				var cert tls.Certificate
				cert, err := tls.LoadX509KeyPair(p.listenerConfigs.Creeds.PublicKey, p.listenerConfigs.Creeds.PrivateKey)
				if err != nil {
					logger.Sugar.Fatal(errLoadTLSKeys(err))
				}
				tlsConfigs = &tls.Config{
					Certificates:       []tls.Certificate{cert},
					InsecureSkipVerify: true,
				}
			}
		}
	}
	listener, err := net.ListenTCP("tcp", p.listenerAddr)
	if err != nil {
		logger.Sugar.Fatalf("Failed to open local port to listen: %s", err)
	}

	p.addClientsConnectionsLoop()
	p.delClientsConnectionsLoop()

	for !p.closed.IsSet() {
		if tlsConfigs != nil {
			var (
				netConn net.Conn
			)
			netConn, err = listener.AcceptTCP()
			logger.Cyan.Printf("New incoming proxy connection: [%s]\n", netConn.LocalAddr().String())
			if err != nil {
				logger.Sugar.Errorf("Failed to accept connection '%s'", err)
				continue
			}
			p.listenerConn = tls.Server(netConn, tlsConfigs)
		} else {
			p.listenerConn, err = listener.AcceptTCP()
			logger.Cyan.Printf("New incoming proxy connection: [%s]\n", p.listenerAddr.String())
			if err != nil {
				logger.Sugar.Errorf("Failed to accept connection '%s'", err)
				continue
			}
		}
		p.setTCPFlags(p.listenerConn, p.listenerConfigs)
		p.setNetFlags(p.listenerConn, p.listenerConfigs)
		//for i := 0; i < runtime.NumCPU(); i++ {
		for j := range addrList {
			buf := *addrList[j]
			if buf.UseKeepAlive && buf.KeepAliveTimeout == 0 {
				buf.KeepAliveTimeout = consts.KeepAliveTimeout
			}
			p.EmmitAddSignal(&models.ConnSignal{
				OperationType: enums.ADD,
				ConnConfigs:   &buf,
			})
		}
		//	}
		logger.Green.Printf("Start proxy, type: [%s], addr: [%s]\n", p.listenerConfigs.Protocol, p.listenerAddr.String())
	}
}
