package tcp_ip_proxy

import "time"

type TCPConnection interface {
	SetNoDelay(bool) error
	SetKeepAlive(keepalive bool) error
	SetKeepAlivePeriod(d time.Duration) error
	SetLinger(sec int) error
}

type Proxy interface{}
