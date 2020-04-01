package main

import (
	"fmt"
	"io"
	"net"

	tcp_ip_proxys "github.com/alex60217101990/proxy.git/external/tcp-ip-proxys"
)

func main() {
	var (
		conn io.ReadWriteCloser
		err  error
	)

	conn, err = net.DialTCP("tcp", &net.TCPAddr{}, &net.TCPAddr{})
	if conn == nil || err != nil {
		fmt.Println(err, conn)
	}

	// conn, err = tls.Dial("tcp", "localhost:8000", nil)
	// if conn == nil || err != nil {
	// 	fmt.Println(err, conn)
	// }
	test(conn)
}

func test(lconn io.ReadWriteCloser) {
	if c, ok := lconn.(tcp_ip_proxys.TCPConnection); lconn != nil && ok {
		fmt.Println(c, "implement")
	}
}
