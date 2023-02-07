package pac

import "github.com/dop251/goja"

// HostExtension allow to extend VM to customize PAC agent
type HostExtension func(vm *goja.Runtime) (*goja.Runtime, error)

// A HostAdapter can be specified native PAC functions when you initialize PAC
// executor. it allowed you to use your own function definition.
type HostAdapter struct {
	DNSResolve  func(host string) string
	MyIPAddress func() string
	Alert       func(msg string)
}

//
func (adapt *HostAdapter) init(parent *HostAdapter) {
	if parent != nil {
		*adapt = *parent
	}
	if adapt.DNSResolve == nil {
		adapt.DNSResolve = defaultDNSResolve
	}
	if adapt.MyIPAddress == nil {
		adapt.MyIPAddress = defaultMyIPAddress
	}
	if adapt.Alert == nil {
		adapt.Alert = defaultAlert
	}
}

//
func (adapt *HostAdapter) apply(vm *goja.Runtime) error {
	if err := vm.Set("dnsResolve", adapt.DNSResolve); err != nil {
		return err
	}
	if err := vm.Set("myIpAddress", adapt.MyIPAddress); err != nil {
		return err
	}
	if err := vm.Set("alert", adapt.Alert); err != nil {
		return err
	}
	return nil
}
