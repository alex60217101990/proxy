package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alex60217101990/proxy.git/external/enums"
	"github.com/alex60217101990/proxy.git/external/logger"
	"github.com/alex60217101990/proxy.git/external/models"
	tcp_ip_proxys "github.com/alex60217101990/proxy.git/external/tcp-ip-proxys"
)

func main() {
	logger.InitLogger()

	testHTTPServer()
	testTCPProxyServer()

	var Stop = make(chan os.Signal)
	signal.Notify(Stop,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGABRT,
	)
	for {
		select {
		case <-Stop:
			logger.Close()
			return
		}
	}
}

func testTCPProxyServer() {
	go func() {
		proxy := tcp_ip_proxys.NewTCPProxy(
			tcp_ip_proxys.SetConfigs(&models.ConnConfigs{
				Protocol: enums.TCP,
				Deadline: 15,
				//MaxBufferSize:    1 << 20,
				UseKeepAlive:     true,
				KeepAliveTimeout: 15,
				LisAddress:       "localhost:2233",
			}))
		defer proxy.Close()
		proxy.Listen([]*models.ConnConfigs{
			&models.ConnConfigs{
				Deadline: 15,
				//MaxBufferSize:    1 << 20,
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
