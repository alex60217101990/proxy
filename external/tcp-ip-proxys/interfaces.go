package tcp_ip_proxys

import (
	"time"

	"github.com/alex60217101990/proxy.git/external/models"
)

type TCPConnection interface {
	SetNoDelay(bool) error
	SetKeepAlive(keepalive bool) error
	SetKeepAlivePeriod(d time.Duration) error
	SetLinger(sec int) error
}

type NetConnection interface {
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

type Proxy interface {
	Close()
	Listen(addrList ...*models.ConnConfigs)
	EmmitAddSignal(signal *models.ConnSignal)
	EmmitDelSignal(signal *models.ConnSignal)
}

type Firewall interface {
	
}