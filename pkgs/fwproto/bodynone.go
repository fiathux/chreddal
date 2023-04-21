package fwproto

import "fmt"

// BodyNone represents the body of the frame with no payload.
type BodyNone struct {
}

// Parse parses the body from the byte slice.
func (bn BodyNone) Parse(b []byte) (int, error) {
	return 0, nil
}

// Serilize serilizes the body to the byte slice.
func (bn BodyNone) Serilize(prevBuffer []byte) ([]byte, error) {
	return prevBuffer, nil
}

// String returns the string representation of the body.
func (bn BodyNone) String() string {
	return fmt.Sprintf("Body(None)")
}
