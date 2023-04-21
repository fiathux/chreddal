// Package socks implements a SOCKS (SOCKS4, SOCKS4A and SOCKS5) proxy server.
//
// The proxy server could set a forwarder to forward the request on different
// way. such as: SSL, SSH, Websocket, etc.
//
// If no forwarder is specified, the proxy server will work as a normal SOCKS
// server.
package socks
