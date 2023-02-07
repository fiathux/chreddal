package pac

import (
	"testing"

	"github.com/dop251/goja"
)

const nativeTestScript = `
function outx(){
	alert("Do output")
 	return [
 		dnsResolve("www.google.com"),
		myIpAddress(),
	]
}

function FindProxyForURL(){}
`

// Native function test
func TestNative(t *testing.T) {
	agt, err := initAgent("testing.js", nativeTestScript, nil)
	if err != nil {
		t.Log("Failed to create VM - ", err)
		t.FailNow()
	}
	vm = agt.(*scriptVM).vm
	fn, ok := goja.AssertFunction(vm.Get("outx"))
	if !ok {
		t.Log("No function to embedded into VM")
		t.FailNow()
	}
	rst, err := fn(goja.Undefined())
	if err != nil {
		t.Log("Failed to run embedded test function:", err)
		t.FailNow()
	}
	obj, ok := rst.Export().([]any)
	if dns, ok := obj[0].(string); ok && dns != "" {
		t.Log("dnsResolve:", dns)
	} else {
		t.Log("failed execute dnsResolve")
		t.Fail()
	}
	if myip, ok := obj[1].(string); ok && myip != "" {
		t.Log("myIpAddress:", myip)
	} else {
		t.Log("failed execute myIpAddress")
		t.Fail()
	}
}
