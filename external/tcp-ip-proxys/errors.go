package tcp_ip_proxys

import "fmt"

var (
	errEmptyProxyListenerConfigs = func(params ...interface{}) error {
		return fmt.Errorf("%s proxy server listener configs is empty (<nil>)", params...)
	}
	errEmptyProxyListenerAddress = func(params ...interface{}) error {
		return fmt.Errorf("%s proxy server listener address string parameter is empty", params...)
	}
	errResolveLocalAddress = func(param interface{}) error {
		return fmt.Errorf("failed to resolve local address: %s", param)
	}
	errToSmallOptionsList = func(param interface{}) error {
		return fmt.Errorf("to small options list: %s", param)
	}
	errParseRootCert = func(params ...interface{}) error {
		return fmt.Errorf("proxy: type: [%v], addr [%s]. Failed to parse root certificate", params...)
	}
	errLoadTLSKeys = func(param interface{}) error {
		return fmt.Errorf("load TLS keys error: %v", param)
	}
)
