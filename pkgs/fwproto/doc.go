// Package fwproto impolments the protocol for data forwarding between each
// proxy server.
//
// The protocol is supported full feature of SOCKS5. including bind and UDP.
//
// It also supports reverse proxy that from public network to private network.
//
// A notification mechanism is implemented to create tunnel as needed. it can
// be use to pass the customized message for user applications as well.
package fwproto
