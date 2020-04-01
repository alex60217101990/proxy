package models

import "net"

type UDPPacket struct {
	Src  *net.UDPAddr
	Data []byte
}
