package fwproto

import "fmt"

// BodyPayloadStream represents the payload stream of the body.
type BodyPayloadStream struct {
	Data []byte
}

// Parse parses the payload stream from the byte slice.
func (bps BodyPayloadStream) Parse(b []byte) (int, error) {
	bps.Data = b
	return len(b), nil
}

// Serilize serilizes the payload stream to the byte slice.
func (bps BodyPayloadStream) Serilize(prevBuffer []byte) ([]byte, error) {
	return append(prevBuffer, bps.Data...), nil
}

// String returns the string representation of the payload stream.
func (bps BodyPayloadStream) String() string {
	return fmt.Sprintf("Payload Stream - %v", bps.Data)
}
