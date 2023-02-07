package pac

import (
	"testing"

	"github.com/dop251/goja"
)

const phonyResolve = `
function dnsResolve(host) {
  return "127.0.0.1"
}
`

var vm = goja.New()

func execJSFunc(
	t *testing.T, name string,
	runner func(goja.Callable) (goja.Value, error),
) {
	f, ok := goja.AssertFunction(vm.Get(name))
	if !ok {
		t.Log(name + " is not a function")
		t.FailNow()
	}
	if r, err := runner(f); err != nil {
		t.Log("Error: "+name+" - ", err)
		t.FailNow()
	} else {
		t.Log(name+":", r)
	}
}

func TestBuildInJSAvailable(t *testing.T) {
	t.Log(vm.RunString(phonyResolve))
	t.Log(vm.RunProgram(buildinScript))
	// isPlainHostName
	execJSFunc(t, "isPlainHostName", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("www"))
	})
	// dnsDomainIs
	execJSFunc(t, "dnsDomainIs", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("www.xyz.com"), vm.ToValue("com"))
	})
	// localHostOrDomainIs
	execJSFunc(t, "localHostOrDomainIs", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("www"), vm.ToValue("www.xyz.com"))
	})
	// dnsDomainLevels
	execJSFunc(t, "dnsDomainLevels", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("www.xyz.com"))
	})
	// shExpMatch
	execJSFunc(t, "shExpMatch", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("www.xyz.com"), vm.ToValue("*.xyz.*"))
	})
	// weekdayRange
	execJSFunc(t, "weekdayRange", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("MON"), vm.ToValue("FRI"))
	})
	// dateRange
	execJSFunc(t, "dateRange", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(),
			vm.ToValue(1), vm.ToValue("JAN"), vm.ToValue(1), vm.ToValue("MAY"))
	})
	// timeRange
	execJSFunc(t, "timeRange", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(),
			vm.ToValue(0), vm.ToValue(0), vm.ToValue(9), vm.ToValue(30),
			vm.ToValue("GMT"))
	})
	// isValidIpAddress
	execJSFunc(t, "isValidIpAddress", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("192.168.0.1"))
	})
	// convert_addr
	execJSFunc(t, "convert_addr", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("10.11.0.1"))
	})
	// isResolvable
	execJSFunc(t, "isResolvable", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("www.xyz.com"))
	})
	// isInNet
	execJSFunc(t, "isInNet", func(f goja.Callable) (goja.Value, error) {
		return f(goja.Undefined(), vm.ToValue("www.xyz.com"),
			vm.ToValue("127.0.0.0"), vm.ToValue("255.255.255.0"))
	})
}
