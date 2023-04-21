package fwproto

import "fmt"

// FrameHead represents the head of the frame.
type FrameHead struct {
	Command FrameCommand
	Length  Int24
}

// Frame represents a frame of the protocol.
type Frame struct {
	Head FrameHead
	Body AnyMatch
}

// Parse parses the frame head from the byte slice.
func (fh FrameHead) Parse(b []byte) (int, error) {
	var err error
	var n int
	fh.Command = FrameCommand(b[0])
	if !fh.Command.IsValid() {
		return 0, fmt.Errorf("invalid frame command: %d", fh.Command)
	}
	n, err = fh.Length.Parse(b[1:])
	if err != nil {
		return 0, err
	}
	return n + 1, nil
}

// Serilize serilizes the frame head to the byte slice.
func (fh FrameHead) Serilize(prevBuffer []byte) ([]byte, error) {
	if !fh.Command.IsValid() {
		return nil, fmt.Errorf("invalid frame command: %d", fh.Command)
	}
	b := append(prevBuffer, byte(fh.Command))
	b, err := fh.Length.Serilize(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Parse parses the frame from the byte slice.
func (f Frame) Parse(b []byte) (int, error) {
	var err error
	var n int
	n, err = f.Head.Parse(b)
	if err != nil {
		return n, err
	}
	f.Body = f.Head.Command.ChooseMatcher()
	n, err = f.Body.Parse(b[n:])
	if err != nil {
		return n, err
	}
	return n, nil
}

// Serilize serilizes the frame to the byte slice.
func (f Frame) Serilize(prevBuffer []byte) ([]byte, error) {
	b, err := f.Head.Serilize(prevBuffer)
	if err != nil {
		return nil, err
	}
	b, err = f.Body.Serilize(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String returns the string representation of the frame.
func (f Frame) String() string {
	cmd := f.Head.Command.String()
	return fmt.Sprintf("Frame [%s]: %s", cmd, f.Body)
}
