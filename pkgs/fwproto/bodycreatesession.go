package fwproto

import (
	"fmt"
	"strings"
)

// SessionTypes represents the type of the session.
type SessionTypes uint8

// All SessionTypes
const (
	SesTypeStream        SessionTypes = 0
	SesTypeBind                       = 1
	SesTypeReverseAccept              = 2
	SesTypeDgram                      = 3
)

// BodyCreateSession represents the body of the frame to create session
type BodyCreateSession struct {
	SessionType SessionTypes
	UUID        UUID
	Params      []Param
}

// IsSessionTypeValid returns true if the session type is valid.
func IsSessionTypeValid(t SessionTypes) bool {
	switch t {
	case SesTypeStream, SesTypeBind, SesTypeReverseAccept, SesTypeDgram:
		return true
	default:
		return false
	}
}

// String returns the string representation of the body.
func (st SessionTypes) String() string {
	switch st {
	case SesTypeStream:
		return "Stream"
	case SesTypeBind:
		return "Bind"
	case SesTypeReverseAccept:
		return "ReverseAccept"
	case SesTypeDgram:
		return "Dgram"
	default:
		return "Unknown"
	}
}

// Parse parses the body from the byte slice.
func (bcs BodyCreateSession) Parse(b []byte) (int, error) {
	bcs.SessionType = SessionTypes(b[0])
	if !IsSessionTypeValid(bcs.SessionType) {
		return 0, fmt.Errorf("Invalid session type: %d", bcs.SessionType)
	}
	copy(bcs.UUID[:], b[1:17])
	countParam := int(b[17])
	bcs.Params = make([]Param, countParam)
	n := 18
	for i := 0; i < countParam; i++ {
		var err error
		n, err = bcs.Params[i].Parse(b[n:])
		if err != nil {
			return 0, err
		}
	}
	return n, nil
}

// Serilize serilizes the body to the byte slice.
func (bcs BodyCreateSession) Serilize(prevBuffer []byte) ([]byte, error) {
	if !IsSessionTypeValid(bcs.SessionType) {
		return nil, fmt.Errorf("Invalid session type: %d", bcs.SessionType)
	}
	if len(bcs.Params) > 4 {
		return nil, fmt.Errorf("Too many parameters: %d", len(bcs.Params))
	}
	b := append(prevBuffer, byte(bcs.SessionType))
	b = append(b, bcs.UUID[:]...)
	b = append(b, byte(len(bcs.Params)))
	for _, p := range bcs.Params {
		var err error
		b, err = p.Serilize(b)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

// String returns the string representation of the body.
func (bcs BodyCreateSession) String() string {
	params := make([]string, len(bcs.Params))
	for i, p := range bcs.Params {
		params[i] = p.String()
	}
	return fmt.Sprintf(
		"Body CreateSession - Type:%s UUID:%s\n  Params:\n    %s",
		bcs.SessionType.String(), bcs.UUID, strings.Join(params, "\n    "),
	)
}
