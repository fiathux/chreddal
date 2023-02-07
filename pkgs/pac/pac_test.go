package pac

import (
	"fmt"
	"testing"
)

func noErr(
	t *testing.T, pac Agent, url, host string,
) func(func(*ProxyResult) error) {
	pxs := pac.FindProxyForURL(url, host)
	if pxs.Err != nil {
		t.Log("Error query", url, "-", pxs.Err)
		t.Fail()
		return func(f func(*ProxyResult) error) {}
	}
	return func(f func(*ProxyResult) error) {
		if err := f(&pxs); err != nil {
			t.Log("Error check", url, "-", err)
			t.Fail()
		} else {
			t.Log("Succeed check", url)
		}
	}
}

// PAC script test
func TestPAC(t *testing.T) {
	pac, err := FromFile("test.pac", nil)
	if err != nil {
		t.Log("Failed create agent", err)
		t.FailNow()
	}
	// check regular proxy rule
	checkregproxy := func(pxs *ProxyResult) error {
		if len(pxs.Items) != 1 {
			return fmt.Errorf("invalid return - Len:%d", len(pxs.Items))
		}
		tgt := pxs.Items[0]
		if tgt.Type != PxTHTTP {
			return fmt.Errorf(
				"error parse proxy type %s - %s", tgt.Type.String(), tgt.String())
		}
		if tgt.Target != "localhost" || tgt.Port != 8080 {
			return fmt.Errorf(
				"error parse proxy host %s", tgt.String())
		}
		t.Log(tgt.String())
		return nil
	}
	noErr(t, pac, "https://www.google.com/123",
		"www.google.com")(checkregproxy)
	noErr(t, pac, "https://404.google.cn/nothing",
		"404.google.cn")(checkregproxy)
	noErr(t, pac, "https://name.who.blogspot.com/chapter0",
		"name.who.blogspot.com")(checkregproxy)
	// test SOCKS parse
	noErr(t, pac, "http://sockstest.com",
		"sockstest.com")(func(pxs *ProxyResult) error {
		if len(pxs.Items) != 3 {
			return fmt.Errorf("invalid return - Len:%d", len(pxs.Items))
		}
		ss := pxs.Items[0]
		s5 := pxs.Items[1]
		s4 := pxs.Items[2]
		if ss.Type != PxTSocks {
			return fmt.Errorf("type socks parsed error - %s", ss.String())
		}
		if s5.Type != PxTSocks5 {
			return fmt.Errorf("type socks5 parsed error - %s", s5.String())
		}
		if s4.Type != PxTSocks4 {
			return fmt.Errorf("type socks4 parsed error - %s", s4.String())
		}
		if s5.Port != 1080 || s4.Port != 1080 {
			return fmt.Errorf("default socks port error -\n  %s\n  %s",
				s5.String(), s4.String())
		}
		t.Log("Succeed check", "- [SOCKS]", ss)
		t.Log("Succeed check", "- [SOCKS5]", s5)
		t.Log("Succeed check", "- [SOCKS4]", s4)
		return nil
	})
	// complex result
	noErr(t, pac, "http://complex.com",
		"complex.com")(func(pxs *ProxyResult) error {
		if len(pxs.Items) != 4 {
			return fmt.Errorf("invalid return - Len:%d", len(pxs.Items))
		}
		ht := pxs.Items[0]
		hts := pxs.Items[1]
		s4 := pxs.Items[2]
		natv := pxs.Items[3]
		if ht.Type != PxTHTTP {
			return fmt.Errorf("type HTTP parsed error - %s", ht.String())
		}
		if hts.Type != PxTHTTPS {
			return fmt.Errorf("type HTTPS parsed error - %s", hts.String())
		}
		if natv.Type != PxTNative {
			return fmt.Errorf("type Native parsed error - %s", natv.String())
		}
		if ht.Port != 80 {
			return fmt.Errorf("default HTTP port error - %s", ht.String())
		}
		if hts.Port != 443 {
			return fmt.Errorf("default HTTPS port error - %s", hts.String())
		}
		t.Log("Succeed check", "- [HTTP]", ht)
		t.Log("Succeed check", "- [HTTPS]", hts)
		t.Log("Succeed check", "- [SOCKS4]", s4)
		t.Log("Succeed check", "- [NATIVE]", natv)
		return nil
	})
	// test error report
	pxerr := pac.FindProxyForURL("http://error.com", "error.com")
	if len(pxerr.Items) != 1 {
		t.Log("invalid return - Len:", len(pxerr.Items))
		t.Fail()
	}
	if pxerr.Err == nil {
		t.Log("Failed error parse - no error found")
		t.Fail()
	} else {
		t.Log("successed report errors - ", pxerr.Err.Error())
	}
	pxnoth := pac.FindProxyForURL("http://nothing.com", "nothing.com")
	if len(pxnoth.Items) != 0 || pxnoth.Err == nil {
		t.Log("failed to parse case nothing", pxnoth)
		t.Fail()
	} else {
		t.Log("successed report nothing - ", pxnoth.Err.Error())
	}
	// test direct
	noErr(t, pac, "http://direct.local",
		"direct.local")(func(pxs *ProxyResult) error {
		if len(pxs.Items) != 1 {
			return fmt.Errorf("invalid return - Len:%d", len(pxs.Items))
		}
		if pxs.Items[0].Type != PxTDirect {
			return fmt.Errorf("error parse DIRECT - %s", pxs.Items[0].String())
		}
		return nil
	})
}
