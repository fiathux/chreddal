// Author: Moi Su
// Date: 2022-11-28

package logger

import (
	"io"
	"log"
	"os"
	"sync/atomic"
)

// LevelMask use for disable some log level
type LevelMask uint32

// Disablable log levels
const (
	MskDbg LevelMask = 1 << iota
	MskIfo
	MskWarn
	MskErr
	MskFatal
)

// Some quick mask set
const (
	MskNone  LevelMask = 0
	MskNoDbg           = MskDbg
	MskNoIfo           = MskDbg | MskIfo
)

// an implement of logger stream
type logStreamer struct {
	io     io.Writer
	prefix string
	flag   int
	mask   *LevelMask
	dbg    *log.Logger
	inf    *log.Logger
	warn   *log.Logger
	err    *log.Logger
	fatal  *log.Logger
}

// a logging object who used a format string
type logFmtStream struct {
	logger FmtOutput
	fmt    string
}

// logger stream with format support
type logStreamWithFormat logStreamer

// Output represent a specific logging object
type Output interface {
	Debug(v ...any)
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}

// FmtOutput represent a logging object that using customer format string
type FmtOutput interface {
	Debug(fmt string, v ...any)
	Info(fmt string, v ...any)
	Warn(fmt string, v ...any)
	Error(fmt string, v ...any)
	Fatal(fmt string, v ...any)
}

// Log is an custom logging object
type Log interface {
	Output
	// get formating outputer
	F() FmtOutput
	// fork for new specific log
	Specific(prefix string) Log
	// fork for new specific log with formatter
	FSpecific(fmt string) Output
	// get raw logger
	Raw(name string) *log.Logger
	// set log level mask
	SetMask(mask LevelMask)
}

func inhertLog(
	writer io.Writer, prefix string, flag int, mask *LevelMask,
) Log {
	if prefix == "" {
		prefix = " "
	} else {
		prefix = " " + prefix + " "
	}
	ret := &logStreamer{
		io:     writer,
		prefix: prefix,
		flag:   flag,
		mask:   mask,
		dbg:    log.New(writer, "DEBUG,"+prefix, flag),
		inf:    log.New(writer, "INFO ,"+prefix, flag),
		warn:   log.New(writer, "WARN ,"+prefix, flag),
		err:    log.New(writer, "ERROR,"+prefix, flag),
		fatal:  log.New(writer, "FATAL,"+prefix, flag),
	}
	return ret
}

// New create a new logging object
func New(writer io.Writer, prefix string, flag int, mask LevelMask) Log {
	return inhertLog(writer, prefix, flag, &mask)
}

// StdLog as a default logger that output to stdout
var StdLog = New(os.Stdout, "", log.LstdFlags|log.LUTC, MskNone)

// ErrLog as a default logger that output to stderr
var ErrLog = New(os.Stderr, "", log.LstdFlags|log.LUTC, MskNone)

// Phony log functions
func phonyPrint(v ...any) {
	return
}
func phonyPrintf(fmt string, v ...any) {
	return
}

// check log functions
func logPrint(in func(v ...any), chk bool) func(v ...any) {
	if chk {
		return in
	}
	return phonyPrint
}
func logPrintf(
	in func(fmt string, v ...any), chk bool,
) func(fmt string, v ...any) {
	if chk {
		return in
	}
	return phonyPrintf
}

func (l *logStreamer) maskout(in func(v ...any), msk LevelMask) func(v ...any) {
	return logPrint(in, msk&LevelMask(atomic.LoadUint32((*uint32)(l.mask))) == 0)
}

//
func (l *logStreamer) Debug(v ...any) {
	l.maskout(l.dbg.Println, MskDbg)(v...)
}

//
func (l *logStreamer) Info(v ...any) {
	l.maskout(l.inf.Println, MskIfo)(v...)
}

//
func (l *logStreamer) Warn(v ...any) {
	l.maskout(l.warn.Println, MskWarn)(v...)
}

//
func (l *logStreamer) Error(v ...any) {
	l.maskout(l.err.Println, MskErr)(v...)
}

//
func (l *logStreamer) Fatal(v ...any) {
	l.maskout(l.fatal.Println, MskFatal)(v...)
}

//
func (l *logStreamer) Specific(prefix string) Log {
	return inhertLog(l.io, prefix, l.flag, l.mask)
}

//
func (l *logStreamer) FSpecific(fmt string) Output {
	ret := &logFmtStream{
		logger: (*logStreamWithFormat)(l),
		fmt:    fmt,
	}
	return ret
}

//
func (l *logStreamer) F() FmtOutput {
	return (*logStreamWithFormat)(l)
}

//
func (l *logStreamer) Raw(name string) *log.Logger {
	switch name {
	case "dbg":
		return l.dbg
	case "inf":
		return l.inf
	case "warn":
		return l.warn
	case "err":
		return l.err
	case "fatal":
		return l.fatal
	default:
		return nil
	}
}

//
func (l *logStreamer) SetMask(mask LevelMask) {
	atomic.StoreUint32((*uint32)(l.mask), uint32(mask))
}

////////////////////////

func (l *logStreamWithFormat) maskout(
	in func(fmt string, v ...any),
	msk LevelMask,
) func(fmt string, v ...any) {
	return logPrintf(in, msk&LevelMask(atomic.LoadUint32((*uint32)(l.mask))) == 0)
}

//
func (l *logStreamWithFormat) Debug(fmt string, v ...any) {
	l.maskout(l.dbg.Printf, MskDbg)(fmt, v...)
}

//
func (l *logStreamWithFormat) Info(fmt string, v ...any) {
	l.maskout(l.inf.Printf, MskIfo)(fmt, v...)
}

//
func (l *logStreamWithFormat) Warn(fmt string, v ...any) {
	l.maskout(l.warn.Printf, MskWarn)(fmt, v...)
}

//
func (l *logStreamWithFormat) Error(fmt string, v ...any) {
	l.maskout(l.err.Printf, MskErr)(fmt, v...)
}

//
func (l *logStreamWithFormat) Fatal(fmt string, v ...any) {
	l.maskout(l.fatal.Printf, MskFatal)(fmt, v...)
}

////////////////////////

//
func (l *logFmtStream) Debug(v ...any) {
	l.logger.Debug(l.fmt, v...)
}

//
func (l *logFmtStream) Info(v ...any) {
	l.logger.Info(l.fmt, v...)
}

//
func (l *logFmtStream) Warn(v ...any) {
	l.logger.Warn(l.fmt, v...)
}

//
func (l *logFmtStream) Error(v ...any) {
	l.logger.Error(l.fmt, v...)
}

//
func (l *logFmtStream) Fatal(v ...any) {
	l.logger.Fatal(l.fmt, v...)
}
