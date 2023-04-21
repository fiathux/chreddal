package fwproto

// FrameCommand represents the command of the frame.
type FrameCommand uint8

const (
	CMDAliveTest     FrameCommand = 0x00
	CMDGeneralAck    FrameCommand = 0x01
	CMDCreateSession FrameCommand = 0x02
	CMDSelectSession FrameCommand = 0x03
	CMDCloseSession  FrameCommand = 0x04
	CMDEventNotify   FrameCommand = 0x05
	CMDPayloadStream FrameCommand = 0x60
	CMDPayloadDgram  FrameCommand = 0x61
)

// String returns the string representation of the command.
func (fc FrameCommand) String() string {
	switch fc {
	case CMDAliveTest:
		return "AliveTest"
	case CMDGeneralAck:
		return "GeneralAck"
	case CMDCreateSession:
		return "CreateSession"
	case CMDSelectSession:
		return "SelectSession"
	case CMDCloseSession:
		return "CloseSession"
	case CMDEventNotify:
		return "EventNotify"
	case CMDPayloadStream:
		return "PayloadStream"
	case CMDPayloadDgram:
		return "PayloadDgram"
	default:
		return "Unknown"
	}
}

// ChooseMatcher returns the matcher of the body of specified command.
func (fc FrameCommand) ChooseMatcher() AnyMatch {
	switch fc {
	case CMDAliveTest:
		return &BodyNone{}
	case CMDGeneralAck:
		return &BodyGeneralAck{}
	case CMDCreateSession:
		return &BodyCreateSession{}
	case CMDSelectSession:
		return &BodySelectSession{}
	case CMDCloseSession:
		return &BodyNone{}
	case CMDEventNotify:
		return &BodyEventNotify{}
	case CMDPayloadStream:
		return &BodyPayloadStream{}
	case CMDPayloadDgram:
		return &BodyPayloadDgram{}
	default:
		return nil
	}
}

// IsValid returns whether the command is valid.
func (fc FrameCommand) IsValid() bool {
	switch fc {
	case CMDAliveTest:
	case CMDGeneralAck:
	case CMDCreateSession:
	case CMDSelectSession:
	case CMDCloseSession:
	case CMDEventNotify:
	case CMDPayloadStream:
	case CMDPayloadDgram:
	default:
		return false
	}
	return true
}
