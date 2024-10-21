package debug

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
)

var Enabled bool

func init() {
	Enabled = load()
	if Enabled {
		fmt.Fprintln(os.Stderr, "---- DEBUG MODE ENABLED ----")
	}
}

func load() bool {
	v, ok := os.LookupEnv("DEBUG")
	if !ok {
		return false
	}
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "1" || v == "true"
}

func Breakpoint(dumpVar any, stop bool) {
	if stop {
		return
	}
	fmt.Fprintln(os.Stderr, "-------------------- Breakpoint --------------------")

	fmt.Fprintln(os.Stderr, "Dumping current stack:")
	debug.PrintStack()
	fmt.Fprintln(os.Stderr, "")

	fmt.Fprintln(os.Stderr, "Dumping variable:")
	spew.Fdump(os.Stderr, dumpVar)
	fmt.Fprintln(os.Stderr, "")

	fmt.Fprintf(os.Stderr, "SIGTRAP ignored: %v\n", signal.Ignored(syscall.SIGTRAP))

	runtime.Breakpoint()
}
