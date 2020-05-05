// Package cli contains the dgo CLI (Command Line Interface)
package cli

import (
	"flag"
	"io"

	"github.com/tada/catch/pio"

	"github.com/tada/catch"
	"github.com/tada/dgo/util"
)

// Command is the interface implemented by all dgo commands. Both global and sub commands.
type Command interface {
	// Do parses the arguments and runs the command
	Do([]string) int

	// Help prints the help for the command
	Help()

	// MissingOption prints the missing option error for the given option and returns 1
	MissingOption(string) int

	// RunWithCatch runs the given function, recovers error panics and reports them. It returns a non zero exit status
	// in case an error was recovered
	RunWithCatch(func() int) int

	// Name returns the name of the command
	Name() string

	// Err returns the error output writer
	Err() io.Writer

	// Out returns the output writer
	Out() io.Writer

	// Verbose returns the verbosity
	Verbose() bool
}

type command struct {
	parent  Command
	flags   *flag.FlagSet
	out     io.Writer
	err     io.Writer
	verbose bool
}

func (h *command) Help() {
	if h.parent != nil {
		pio.WriteString(h.out, h.parent.Name())
		pio.WriteByte(h.out, ' ')
	}
	pio.WriteString(h.out, h.Name())
	pio.WriteByte(h.out, '\n')
	h.flags.SetOutput(h.out)
	h.flags.PrintDefaults()
}

func (h *command) Err() io.Writer {
	return h.err
}

func (h *command) Out() io.Writer {
	return h.out
}

func (h *command) Verbose() bool {
	return h.verbose
}

func (h *command) Name() string {
	return h.flags.Name()
}

func (h *command) RunWithCatch(runner func() int) int {
	exitCode := 0
	err := catch.Do(func() {
		exitCode = runner()
	})
	if err != nil {
		util.Fprintf(h.err, "Error: %s\n", err.Error())
		return 1
	}
	return exitCode
}

func (h *command) MissingOption(opt string) int {
	util.Fprintf(h.err, "missing required option: -%s\n", opt)
	return 1
}

func (h *command) Parse(args []string) (int, bool) {
	h.flags.SetOutput(h.err)
	err := h.flags.Parse(args)
	if err != nil {
		if err == flag.ErrHelp {
			h.Help()
			return 0, true
		}
		return 1, true
	}
	return 0, false
}
