package log

import (
	"fmt"
	"io"
	"os"
)

const (
	Silent     = 0
	ErrorLevel = 1
	WarnLevel  = 2
	InfoLevel  = 3
	DebugLevel = 4
)

var Level = WarnLevel

var (
	debugLogger printer
	warnLogger  printer
	infoLogger  printer
	errorLogger printer
)

var defaultTimeFormat = "2006-01-02 15:04:05.000"

func init() {
	initLoggers(os.Stderr, simpleClock{}, defaultTimeFormat)
}

// log initialization is done in its own function so that it can be manipulated from tests
func initLoggers(output io.Writer, clock clock, timeFormat string) {
	debugLogger = newPrinter("DEBUG ", output, clock, timeFormat)
	infoLogger = newPrinter("INFO  ", output, clock, timeFormat)
	warnLogger = newPrinter("WARN  ", output, clock, timeFormat)
	errorLogger = newPrinter("ERROR ", output, clock, timeFormat)
}

type printer struct {
	prefix     string
	output     io.Writer
	clock      clock
	timeFormat string
}

func newPrinter(prefix string, output io.Writer, clock clock, timeFormat string) printer {
	return printer{
		prefix:     prefix,
		output:     output,
		clock:      clock,
		timeFormat: timeFormat,
	}
}

func (p printer) Printf(format string, v ...any) error {
	message := fmt.Sprintf(format, v...)
	return p.Print(message)
}

func (p printer) Print(message string) error {
	formattedTime := p.clock.now().Format(p.timeFormat)
	length := len(message) + len(p.prefix) + len(formattedTime) + 3
	data := make([]byte, length)

	idx := 0

	// prefix
	copy(data[idx:idx+len(p.prefix)], p.prefix)
	idx += len(p.prefix)

	// first space
	data[idx] = ' '
	idx++

	// date
	copy(data[idx:idx+len(formattedTime)], formattedTime)
	idx += len(formattedTime)

	// second space
	data[idx] = ' '
	idx++

	// message
	copy(data[idx:idx+len(message)], message)
	idx += len(message)

	// new line
	data[idx] = '\n'
	idx++

	_, err := p.output.Write(data)
	return err
}

func Debug(message string) {
	if Level >= DebugLevel {
		debugLogger.Print(message)
	}
}

func Info(message string) {
	if Level >= InfoLevel {
		infoLogger.Print(message)
	}
}

func Warn(message string) {
	if Level >= WarnLevel {
		warnLogger.Print(message)
	}
}

func Error(message string) {
	if Level >= ErrorLevel {
		errorLogger.Print(message)
	}
}

func Debugf(format string, v ...any) {
	if Level >= DebugLevel {
		debugLogger.Printf(format, v...)
	}
}

func Infof(format string, v ...any) {
	if Level >= InfoLevel {
		infoLogger.Printf(format, v...)
	}
}

func Warnf(format string, v ...any) {
	if Level >= WarnLevel {
		warnLogger.Printf(format, v...)
	}
}

func Errorf(format string, v ...any) {
	if Level >= ErrorLevel {
		errorLogger.Printf(format, v...)
	}
}
