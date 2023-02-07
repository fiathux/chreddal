package pac

import (
	"chreddal/pkgs/logger"
	"context"
	"encoding/base64"
	"net"
	"strings"
	"time"

	"github.com/dop251/goja"
)

//go:generate /bin/bash gen-script.sh

const dnsTimeout = 2 * time.Second

// load embedded script
func mustLoadScript(b string) *goja.Program {
	bin, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		panic(err)
	}
	return goja.MustCompile("__builtin__.js", string(bin), true)
}

// default DNS resolve
func defaultDNSResolve(host string) string {
	ctx, _ := context.WithTimeout(context.Background(), dnsTimeout)
	addr, err := net.DefaultResolver.LookupHost(ctx, host)
	if err != nil || len(addr) == 0 {
		return ""
	}
	return addr[0]
}

// default method to get local IP address
func defaultMyIPAddress() string {
	iaddrs, _ := net.InterfaceAddrs()
	loopback := "127.0.0.1"
	for _, addr := range iaddrs {
		if addx, ok := addr.(*net.IPNet); ok {
			ip := addx.IP
			ipstr := ip.String()
			if strings.Index(ipstr, ":") != -1 {
				continue
			}
			if ip.IsLoopback() {
				loopback = ipstr
			} else {
				return ipstr
			}
		}
	}
	return loopback
}

var defaultJSLog = logger.StdLog.Specific("PAC alert")

// default alert function for PAC script
func defaultAlert(msg string) {
	defaultJSLog.Info(msg)
}
