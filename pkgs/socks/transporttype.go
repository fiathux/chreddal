package socks

// TransportType represents the transport type of the request
type TransportType int

// available transport type in socks5
const (
	TrsCONNECT = iota
	TrsBIND
	TrsUDP
)
