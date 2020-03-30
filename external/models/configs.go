package models

import "github.com/alex60217101990/proxy.git/external/enums"

type SingleProxyAddrConfigs struct {
	Protocol               enums.ProtocolType
	TLSUnwrapp             bool
	TLSAddress             string
	LisAddress, RecAddress string
	Conn                   *ConnConfigs
}

type ConnConfigs struct {
	Timeout uint
	// for udp connections
	MaxBufferSize uint
	// for tcp connections
	UseNagles        bool
	UseKeepAlive     bool
	KeepAliveTimeout uint
}
