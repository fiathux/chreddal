package fwproto

import (
	"errors"
	"fmt"
	"net"
)

// Addr represents a network address.
type Addr struct {
	Port int
	IP   []byte
}

// Parse parses the address from the byte slice.
func (a Addr) Parse(b []byte) (int, error) {
	a.Port = int(b[0])<<8 + int(b[1])
	if a.Port == 0 {
		return 0, errors.New("invalid parse address - port can't be 0")
	}
	ipLen := int(b[2])
	if ipLen != 4 && ipLen != 16 {
		return 0, errors.New("invalid parse address - ip length must be 4 or 16")
	}
	a.IP = b[3 : 3+ipLen]
	return 3 + ipLen, nil
}

// Serilize serilizes the address to the byte slice.
func (a Addr) Serilize(prevBuffer []byte) ([]byte, error) {
	if a.Port == 0 {
		return nil, errors.New("invalid serilize address - port can't be 0")
	}
	if len(a.IP) != 4 && len(a.IP) != 16 {
		return nil, errors.New("invalid serilize address - ip length must be 4 or 16")
	}
	buffer := make([]byte, 3+len(a.IP))
	buffer[0] = byte(a.Port >> 8)
	buffer[1] = byte(a.Port)
	buffer[2] = byte(len(a.IP))
	copy(buffer[3:], a.IP)
	return append(prevBuffer, buffer...), nil
}

// String returns the string representation of the address.
func (a Addr) String() string {
	ipStr := net.IP(a.IP).String()
	return fmt.Sprintf("%s:%d", ipStr, a.Port)
}

// ----------------------------

// AnyMatch represent any type of the protocol.
type AnyMatch interface {
	Parse([]byte) (int, error)
	Serilize([]byte) ([]byte, error)
	String() string
}

// ----------------------------

type UInt struct {
	Value uint
}

// String returns the string representation of the unsigned integer.
func (u UInt) String() string {
	return fmt.Sprintf("%d", u.Value)
}

// Int16 represents a 16-bit unsigned integer.
type Int16 struct {
	UInt
}

// Parse parses the integer from the byte slice.
func (i Int16) Parse(b []byte) (int, error) {
	i.Value = uint(b[0])<<8 + uint(b[1])
	return 2, nil
}

// Serilize serilizes the integer to the byte slice.
func (i Int16) Serilize(prevBuffer []byte) ([]byte, error) {
	return append(prevBuffer, byte(i.Value>>8), byte(i.Value)), nil
}

// Int24 represents a 24-bit unsigned integer.
type Int24 struct {
	UInt
}

// Parse parses the integer from the byte slice.
func (i Int24) Parse(b []byte) (int, error) {
	i.Value = uint(b[0])<<16 + uint(b[1])<<8 + uint(b[2])
	return 3, nil
}

// Serilize serilizes the integer to the byte slice.
func (i Int24) Serilize(prevBuffer []byte) ([]byte, error) {
	return append(prevBuffer,
		byte(i.Value>>16), byte(i.Value>>8), byte(i.Value)), nil
}

// Int32 represents a 32-bit unsigned integer.
type Int32 struct {
	UInt
}

// Parse parses the integer from the byte slice.
func (i Int32) Parse(b []byte) (int, error) {
	i.Value = uint(b[0])<<24 + uint(b[1])<<16 + uint(b[2])<<8 + uint(b[3])
	return 4, nil
}

// Serilize serilizes the integer to the byte slice.
func (i Int32) Serilize(prevBuffer []byte) ([]byte, error) {
	return append(prevBuffer,
		byte(i.Value>>24), byte(i.Value>>16), byte(i.Value>>8), byte(i.Value)), nil
}

// ----------------------------

// String represents a string.
type String string

// Parse parses the string from the byte slice.
func (s String) Parse(b []byte) (int, error) {
	lenStr := int(b[0])<<8 + int(b[1])
	if lenStr == 0 {
		return 2, nil
	}
	s = String(b[2 : 2+lenStr])
	return 2 + lenStr, nil
}

// Serilize serilizes the string to the byte slice.
func (s String) Serilize(prevBuffer []byte) ([]byte, error) {
	lenStr := len(s)
	if lenStr == 0 {
		return append(prevBuffer, 0, 0), nil
	}
	buffer := make([]byte, 2+lenStr)
	buffer[0] = byte(lenStr >> 8)
	buffer[1] = byte(lenStr)
	copy(buffer[2:], s)
	return append(prevBuffer, buffer...), nil
}

// String returns the string representation of the string.
func (s String) String() string {
	return string(s)
}

// ----------------------------

// UUID represents a UUID.
type UUID [16]byte

// Parse parses the UUID from the byte slice.
func (u UUID) Parse(b []byte) (int, error) {
	copy(u[:], b[:16])
	return 16, nil
}

// Serilize serilizes the UUID to the byte slice.
func (u UUID) Serilize(prevBuffer []byte) ([]byte, error) {
	return append(prevBuffer, u[:]...), nil
}

// String returns the string representation of the UUID.
func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
