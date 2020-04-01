package models

import (
	"io"
	"net"

	"github.com/alex60217101990/proxy.git/external/enums"
)

type PipeConnTCP struct {
	io.PipeReader
	io.PipeWriter
	conn *net.TCPConn
}

type ConnSignal struct {
	OperationType enums.OperationType
	ConnConfigs   *ConnConfigs
}
