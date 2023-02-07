package logger

import "testing"

func TestLog(t *testing.T) {
	StdLog.Dbg("a debug log", "arg2", 1, 2, 3)
	StdLog.Ifo("a info log")
	StdLog.Warn("a warning log")
	StdLog.Err("a error log")
	StdLog.Fatal("a fatal log")
	erx := ErrLog.Specific("XYZ")
	erx.SetMask(MskNoIfo)
	ErrLog.Dbg("stderr - a debug log")
	ErrLog.Ifo("stderr - a info log")
	ErrLog.Warn("stderr - a warning log")
	ErrLog.Err("stderr - a error log")
	ErrLog.Fatal("stderr - a fatal log")
	erx.Dbg("erx - debug")
	erx.Err("erx - error")
	spx := StdLog.Specific("SPX")
	spx.Dbg("debug")
	spx.F().Ifo("%d - %s", 10, "xyz")
	spxf := spx.FSpecific("%s, %s")
	spxf.Warn("W", "msg")
	spxf.Ifo("I", "x msg")
}
