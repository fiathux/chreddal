package pac

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dop251/goja"
)

// Agent represent a PAC execute Agent
type Agent interface {
	FindProxyForURL(url, host string) ProxyResult
	GetRaw() []byte
}

// PAC VM implemention
type scriptVM struct {
	vm        *goja.Runtime
	adapt     HostAdapter
	rawScript string
	fn        goja.Callable
}

// init PAC agent
func initAgent(
	scrName string, scrContent string,
	adapt *HostAdapter,
	exts ...HostExtension,
) (Agent, error) {
	vm := &scriptVM{
		vm: goja.New(),
	}
	vm.adapt.init(adapt)
	if err := vm.adapt.apply(vm.vm); err != nil {
		return nil, err
	}
	// parse PAC script
	if _, err := vm.vm.RunProgram(buildinScript); err != nil {
		panic("Builtin script error - " + err.Error())
	}
	vm.rawScript = scrContent
	if _, err := vm.vm.RunScript(scrName, scrContent); err != nil {
		return nil, err
	}
	// get entry
	fn, ok := goja.AssertFunction(vm.vm.Get("FindProxyForURL"))
	if !ok {
		return nil, fmt.Errorf("No function named 'FindProxyForURL' in PAC")
	}
	vm.fn = fn
	// setup extensions
	if len(exts) > 0 {
		for _, ex := range exts {
			exvm, err := ex(vm.vm)
			if exvm != nil {
				vm.vm = exvm
			}
			if err != nil {
				return nil, err
			}
		}
	}
	return vm, nil
}

// check to create proxy rule
func makeProxyObj(t ProxyType, host string, dftport int) (ProxyNext, error) {
	if host == "" {
		return ProxyNext{}, fmt.Errorf("Invalid host(0): %s", host)
	}
	// try to split host and port
	n := strings.LastIndex(host, ":")
	xn := strings.Index(host, ":")
	var h, ps string
	if n == 1 {
		return ProxyNext{}, fmt.Errorf("Invalid host(1): %s", host)
	} else if n < 1 || (n != xn && (host[0] != '[' || host[n-1] != ']')) {
		h = host
	} else {
		h = host[:n]
		ps = host[n+1:]
	}
	if ps == "" {
		return ProxyNext{
			Type:   t,
			Target: h,
			Port:   dftport,
		}, nil
	}
	// parse port
	p, err := strconv.Atoi(ps)
	if err != nil || p == 0 || p > 65535 {
		return ProxyNext{}, fmt.Errorf("Invalid host, port error: %s", host)
	}
	return ProxyNext{
		Type:   t,
		Target: h,
		Port:   p,
	}, nil
}

// parse a proxy rule
func parseProxyRule(rule string) (rtype, host string) {
	ru := strings.TrimSpace(rule)
	if ru == "" {
		return "", ""
	}
	n := strings.Index(ru, " ")
	var t, h string
	if n > 0 {
		t = ru[:n]
		h = strings.TrimSpace(ru[n:])
	} else {
		t = ru
	}
	return t, h
}

// FindProxyForURL is port to PAC entry.
//
// Both proxy rules and error might be returned. PAC cloud given several rules.
// so when some rule can not be parsed then error returned, other rules which
// can be parse succeed will be returned yet
func (vm *scriptVM) FindProxyForURL(url, host string) ProxyResult {
	rst, err := vm.fn(goja.Undefined(), vm.vm.ToValue(url), vm.vm.ToValue(host))
	if err != nil {
		return ProxyResult{Raw: "", Err: err}
	}
	rulestr, ok := rst.Export().(string)
	if !ok || rulestr == "" {
		return ProxyResult{Raw: "", Err: fmt.Errorf("no rules")}
	}
	sprule := strings.Split(rulestr, ";")
	ret := make([]ProxyNext, 0, len(sprule))
	var rerr ruleErr
	for _, ru := range sprule {
		t, h := parseProxyRule(ru)
		if t == "" {
			rerr = rerrAppend(rerr, fmt.Sprintf("Invalid proxy rule - %s", ru))
			continue
		}
		var pobj ProxyNext
		var err error
		switch t {
		case "DIRECT":
			pobj = ProxyNext{
				Type:   PxTDirect,
				Target: "",
				Port:   0,
			}
		case "NATIVE":
			pobj = ProxyNext{
				Type:   PxTNative,
				Target: "",
				Port:   0,
			}
		case "HTTP", "PROXY":
			pobj, err = makeProxyObj(PxTHTTP, h, 80)
		case "HTTPS":
			pobj, err = makeProxyObj(PxTHTTPS, h, 443)
		case "SOCKS":
			pobj, err = makeProxyObj(PxTSocks, h, 1080)
		case "SOCKS4":
			pobj, err = makeProxyObj(PxTSocks4, h, 1080)
		case "SOCKS5":
			pobj, err = makeProxyObj(PxTSocks5, h, 1080)
		default:
			rerr = rerrAppend(rerr, fmt.Sprintf("Invalid proxy rule - %s", ru))
			continue
		}
		if err != nil {
			rerr = rerrAppend(rerr,
				fmt.Sprintf("Invalid proxy rule, %s - %s", err.Error(), ru))
			continue
		}
		ret = append(ret, pobj)
	}
	if len(ret) == 0 && rerr == nil {
		return ProxyResult{Raw: rulestr, Err: fmt.Errorf("no rules")}
	}
	return ProxyResult{
		Raw:   rulestr,
		Items: ret,
		Err:   rerr.Export(),
	}
}

// GetRaw get raw PAC script
func (vm *scriptVM) GetRaw() []byte {
	if len(vm.rawScript) == 0 {
		return nil
	}
	ret := make([]byte, len(vm.rawScript))
	copy(ret, vm.rawScript)
	return ret
}
