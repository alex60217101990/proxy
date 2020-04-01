package models

import "github.com/alex60217101990/proxy.git/external/enums"

type SingleProxyAddrConfigs struct {
	IsDebug bool
	// TLSUnwrapp             bool
	// TLSAddress             string
	// LisAddress, RecAddress string
	Conn *ConnConfigs
}

type ConnConfigs struct {
	Protocol    enums.ProtocolType
	Concurrency uint8
	Deadline    uint
	// for udp connections
	MaxBufferSize uint
	// for tcp connections
	UseNagles, UseKeepAlive bool
	LisAddress, RecAddress  string
	KeepAliveTimeout        uint
	LingerSec               int
	Retry                   *Retry
	Creeds                  *Credentials
}

type Retry struct {
	AttemptsNumber  uint8
	AttemptInterval uint
}

type Credentials struct {
	PublicKey  string
	PrivateKey string
	RootCAs    string
}
