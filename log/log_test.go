package log

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testClock struct {
	time time.Time
}

func (c testClock) now() time.Time {
	return c.time
}

func Test(t *testing.T) {
	assert.Equal(t, Level, WarnLevel)

	out := bytes.Buffer{}
	clock := testClock{time: mustParseTime("2006-01-02", "2024-01-01")}
	timeFormat := "2006-01-02 15:04:05.000"
	initLoggers(&out, &clock, timeFormat)

	Level = DebugLevel

	out.Reset()
	clock.time = clock.time.Add(1 * time.Hour)
	Debugf("message %d", 1)
	assert.Equal(t, "DEBUG  2024-01-01 01:00:00.000 message 1\n", out.String())

	out.Reset()
	clock.time = clock.time.Add(1 * time.Hour)
	Infof("message %d", 2)
	assert.Equal(t, "INFO   2024-01-01 02:00:00.000 message 2\n", out.String())

	out.Reset()
	clock.time = clock.time.Add(1 * time.Hour)
	Warnf("message %d", 3)
	assert.Equal(t, "WARN   2024-01-01 03:00:00.000 message 3\n", out.String())

	out.Reset()
	clock.time = clock.time.Add(1 * time.Hour)
	Errorf("message %d", 4)
	assert.Equal(t, "ERROR  2024-01-01 04:00:00.000 message 4\n", out.String())
}

func mustParseTime(format string, s string) time.Time {
	tm, err := time.Parse(format, s)
	if err != nil {
		panic(err)
	}
	return tm
}
