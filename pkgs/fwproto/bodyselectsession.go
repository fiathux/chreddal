package fwproto

import "fmt"

type BodySelectSession struct {
	UUID UUID
}

func (bss BodySelectSession) Parse(b []byte) (int, error) {
	var err error
	var n int
	n, err = bss.UUID.Parse(b)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (bss BodySelectSession) Serilize(prevBuffer []byte) ([]byte, error) {
	b, err := bss.UUID.Serilize(prevBuffer)
	if err != nil {
		return nil, err
	}
	return append(prevBuffer, b...), nil
}

func (bss BodySelectSession) String() string {
	return fmt.Sprintf("Body SelectSession - %v", bss.UUID)
}
