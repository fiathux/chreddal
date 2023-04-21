package fwnet

import (
	"context"
)

// ForwardAddrV4 is a IPv4 address
type ForwardAddrV4 struct {
	IP   [4]byte
	Port uint16
}

// ForwardAddrV6 is a IPv6 address
type ForwardAddrV6 struct {
	IP   [16]byte
	Port uint16
}

// ForwardAddress is a address that compatible with IPv4 and IPv6
type ForwardAddress interface {
	ForwardAddrV4 | ForwardAddrV6
}

// Forwarder is a interface that can create a forward session
type Forwarder[ADDR ForwardAddress] interface {
	Forward(ctx context.Context, dist ADDR) (ForwardSession, error)
	Bind(ctx context.Context, dist ADDR) (BindHandler, error)
	UDPAssociate(ctx context.Context, dist ADDR) (UDPSession[ADDR], error)
}

// ForwardSession is stream session that do proxy forward
type ForwardSession interface {
	Read(data []byte) (int, error)
	Write(buffer []byte) (int, error)
}

// BindHandler is use to implement a bind handler
type BindHandler func() (ForwardSession, error)

// UDPSession is a udp session that do proxy forward
type UDPSession[ADDR ForwardAddress] interface {
	ReadFrom(data []byte) (int, ADDR, error)
	WriteTo(buffer []byte, src ADDR) (int, error)
}
