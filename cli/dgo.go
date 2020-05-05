package cli

import (
	"flag"
	"io"

	"github.com/tada/dgo/util"

	"github.com/tada/catch/pio"
)

// Dgo creates the global dgo command
func Dgo(out, err io.Writer) Command {
	c := &dgoCommand{command: command{out: out, err: err, verbose: false}}
	flags := flag.NewFlagSet(`dgo`, flag.ContinueOnError)
	flags.BoolVar(&c.verbose, `verbose`, false, `Be verbose in output`)
	flags.Usage = c.Help
	c.flags = flags
	return c
}

type dgoCommand struct {
	command
}

func (h *dgoCommand) Do(args []string) int {
	return h.RunWithCatch(func() int {
		r, done := h.Parse(args)
		if done {
			return r
		}
		args = h.flags.Args()
		if len(args) == 0 {
			util.Fprintf(h.err, "missing required command\n")
			return 1
		}
		switch args[0] {
		case `help`:
			if len(args) > 1 {
				switch args[1] {
				case `validate`:
					Validate(h).Help()
				case `help`:
					pio.WriteString(h.out, `prints the help text`)
				default:
					util.Fprintf(h.err, `unknown command: %s`, args[0])
				}
			} else {
				h.Help()
			}
		case `validate`:
			r = Validate(h).Do(args[1:])
		default:
			util.Fprintf(h.err, `unknown command: %s`, args[0])
			r = 1
		}
		return r
	})
}

func (h *dgoCommand) Help() {
	pio.WriteString(h.out, `dgo: a command line tool to interact with the dgo type system

Usage: 
  dgo [flags] <command> [command flags]

Available commands:
  help        Shows this help
  validate    Validates input file against a parameter description

Available flags:
  -verbose   Be verbose in output

Use "dgo help <command>" for more information about a command.
`)
}
