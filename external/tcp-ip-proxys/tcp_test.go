package tcp_ip_proxys

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/alex60217101990/proxy.git/external/enums"
	"github.com/alex60217101990/proxy.git/external/models"
)

func ImplementTest(t *testing.T) {
	var (
		conn io.ReadWriteCloser
		err  error
	)
	t.Run("first", func(t *testing.T) {
		conn, err = net.DialTCP("tcp", &net.TCPAddr{}, &net.TCPAddr{})
		if conn == nil || err != nil {
			t.Error(err, conn)
			return
		}
		if test(conn) {
			t.Log("implement success")
			return
		}
		t.Error("dasn't implement")
	})
	t.Run("second", func(t *testing.T) {
		conn, err = tls.Dial("tcp", "localhost:8000", nil)
		if conn == nil || err != nil {
			t.Error(err, conn)
			return
		}
		if test(conn) {
			t.Log("implement success")
			return
		}
		t.Error("dasn't implement")
	})
}

func testTCPProxyServer() {
	go func() {
		proxy := NewTCPProxy(SetConfigs(&models.ConnConfigs{
			Protocol:         enums.TCP,
			Deadline:         15,
			MaxBufferSize:    1 << 20,
			UseKeepAlive:     true,
			KeepAliveTimeout: 15,
			LisAddress:       "localhost",
		}))
		defer proxy.Close()
		proxy.Listen([]*models.ConnConfigs{
			&models.ConnConfigs{
				Deadline:         15,
				MaxBufferSize:    1 << 20,
				UseKeepAlive:     true,
				KeepAliveTimeout: 15,
				LisAddress:       "localhost:50080",
			},
		}...)
	}()
}

func testHTTPServer() {
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "OK!")
		})
		http.ListenAndServe(":50080", nil)
	}()
}

func test(lconn io.ReadWriteCloser) bool {
	if c, ok := lconn.(TCPConnection); lconn != nil && ok {
		fmt.Println(c, "implement")
		return true
	}
	return false
}

func TCPProxy_TestHTTP(t *testing.T) {
	fmt.Println("111")
	testHTTPServer()
	testTCPProxyServer()
}
