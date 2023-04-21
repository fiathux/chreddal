package fwproto

import "fmt"

// BodyGeneralAck represents the body of the frame with general ack payload.
type BodyGeneralAck struct {
	Result Int16
}

// Parse parses the body from the byte slice.
func (bgk BodyGeneralAck) Parse(b []byte) (int, error) {
	var err error
	var n int
	n, err = bgk.Result.Parse(b)
	if err != nil {
		return n, err
	}
	return n, nil
}

// Serilize serilizes the body to the byte slice.
func (bgk BodyGeneralAck) Serilize(prevBuffer []byte) ([]byte, error) {
	b, err := bgk.Result.Serilize(prevBuffer)
	if err != nil {
		return nil, err
	}
	return append(prevBuffer, b...), nil
}

// String returns the string representation of the body.
func (bgk BodyGeneralAck) String() string {
	return fmt.Sprintf("Body GeneralAck - %v", bgk.Result)
}
