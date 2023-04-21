package fwproto

import "fmt"

type ParamType uint8

const (
	PTTargetAddress ParamType = 0x00
	PTIdleTimeout   ParamType = 0x01
	PTAliveTimeout  ParamType = 0x02
)

// String returns the string representation of the parameter type.
func (pt ParamType) String() string {
	switch pt {
	case PTTargetAddress:
		return "TargetAddress"
	case PTIdleTimeout:
		return "IdleTimeout"
	case PTAliveTimeout:
		return "AliveTimeout"
	default:
		return "Unknown"
	}
}

// IsValid returns whether the parameter type is valid.
func (pt ParamType) IsValid() bool {
	switch pt {
	case PTTargetAddress, PTIdleTimeout, PTAliveTimeout:
		return true
	default:
		return false
	}
}

// ChooseMatcher returns the matcher of the value of specified parameter type.
func (pt ParamType) ChooseMatcher() AnyMatch {
	switch pt {
	case PTTargetAddress:
		return &Addr{}
	case PTIdleTimeout:
		return &Int32{}
	case PTAliveTimeout:
		return &Int32{}
	default:
		return nil
	}
}

// Param represents a parameter in the protocol.
type Param struct {
	Type  ParamType
	Value AnyMatch
}

// Parse parses the parameter from the byte slice.
func (p Param) Parse(b []byte) (int, error) {
	p.Type = ParamType(b[0])
	if !p.Type.IsValid() {
		return 0, fmt.Errorf("invalid parameter type: %d", p.Type)
	}
	p.Value = p.Type.ChooseMatcher()
	n, err := p.Value.Parse(b[1:])
	if err != nil {
		return 0, err
	}
	return n + 1, nil
}

// Serilize serilizes the parameter to the byte slice.
func (p Param) Serilize(prevBuffer []byte) ([]byte, error) {
	if !p.Type.IsValid() {
		return nil, fmt.Errorf("invalid parameter type: %d", p.Type)
	}
	b := append(prevBuffer, byte(p.Type))
	b, err := p.Value.Serilize(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String returns the string representation of the parameter.
func (p Param) String() string {
	return fmt.Sprintf("Param [%s]: %s", p.Type, p.Value)
}
