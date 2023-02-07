package pac

import "fmt"

// ProxyType represent type of proxy. current supported below:
//    PxTDirect		Direct connect(no proxy)
//    PxTHTTP			Use HTTP proxy
//    PxTHTTPS		Use HTTP proxy with TLS connect
//    PxTSocks		Use socks proxy(auto detect socks version)
//    PxTSocks4		Use socks4 proxy
//		PxTSocks5		Use socks5 proxy
//    PxTNative		Use native proxy (non-stdandard)
type ProxyType int

// Allowed proxy type
const (
	PxTDirect = iota + 1
	PxTHTTP
	PxTHTTPS
	PxTSocks
	PxTSocks4
	PxTSocks5
	PxTNative
)

// ProxyNext represent a rule of PAC result
type ProxyNext struct {
	Type   ProxyType // proxy type
	Target string    // target hostname
	Port   int       // host port
}

// ProxyResult represent results of PAC found
type ProxyResult struct {
	Raw   string      // Raw result
	Items []ProxyNext // parsed reuslt
	Err   error       // error messages
}

// retrieve type name
func (t ProxyType) String() string {
	switch t {
	case PxTDirect:
		return "DIRECT"
	case PxTHTTP:
		return "HTTP"
	case PxTHTTPS:
		return "HTTPS"
	case PxTSocks:
		return "SOCKS"
	case PxTSocks4:
		return "SOCKS4"
	case PxTSocks5:
		return "SOCKS5"
	case PxTNative:
		return "NATIVE"
	default:
		return "UNKNOW"
	}
}

// convert ProxyNext to string
func (p *ProxyNext) String() string {
	if p == nil {
		return "NONE"
	}
	if p.Target != "" {
		return fmt.Sprintf("%s %s:%d", p.Type.String(), p.Target, p.Port)
	}
	return p.Type.String()
}
