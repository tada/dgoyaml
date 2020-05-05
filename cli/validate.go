package cli

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tada/catch"
	"github.com/tada/catch/pio"
	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgo/util"
	"github.com/tada/dgoyaml/yaml"
)

// Validate is the Dgo sub command that reads and validates YAML input using a spec written in
// YAML or Dgo
func Validate(parent Command) Command {
	vc := &validateCommand{command: command{parent: parent, out: parent.Out(), err: parent.Err(), verbose: parent.Verbose()}}
	flags := flag.NewFlagSet(`validate`, flag.ContinueOnError)
	flags.StringVar(&vc.input, `input`, ``, `yaml file containing input to validate`)
	flags.StringVar(&vc.spec, `spec`, ``, `yaml or dgo file with the parameter definitions`)
	vc.flags = flags
	return vc
}

type validateCommand struct {
	command
	input string
	spec  string
}

func readFileOrPanic(name string) []byte {
	/* #nosec */
	bs, err := ioutil.ReadFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			err = catch.Error(err)
		}
		panic(err)
	}
	return bs
}

func (h *validateCommand) run() int {
	iMap := h.loadParameters(h.input)
	sType := h.loadStructMapType(h.spec)
	ok := true
	if h.verbose {
		bld := util.NewIndenter(`  `)
		ok = sType.(dgo.MapValidation).ValidateVerbose(iMap, bld)
		pio.WriteString(h.out, bld.String())
	} else {
		vs := sType.(dgo.MapValidation).Validate(nil, iMap)
		if len(vs) > 0 {
			ok = false
			for _, err := range vs {
				pio.WriteString(h.out, err.Error())
				pio.WriteRune(h.out, '\n')
			}
		}
	}
	if ok {
		return 0
	}
	return 1
}

func (h *validateCommand) loadParameters(input string) (iMap dgo.Map) {
	switch {
	case strings.HasSuffix(input, `.yaml`), strings.HasSuffix(input, `.json`):
		data := readFileOrPanic(input)
		m, err := yaml.Unmarshal(data)
		if err != nil {
			panic(catch.Error(err))
		}
		var ok bool
		iMap, ok = m.(dgo.Map)
		if !ok {
			panic(catch.Error(`expecting data to be a map`))
		}
		if h.verbose {
			bld := util.NewIndenter(`  `)
			bld.Append(`Got input yaml with:`)
			b2 := bld.Indent()
			b2.NewLine()
			b2.AppendIndented(string(data))
			pio.WriteString(h.out, bld.String())
		}
	default:
		panic(catch.Error(`invalid file name '%s', expected file name to end with .yaml or .json`, input))
	}
	return
}

func (h *validateCommand) loadStructMapType(spec string) (sType dgo.StructMapType) {
	switch {
	case strings.HasSuffix(spec, `.yaml`), strings.HasSuffix(spec, `.json`):
		m, err := yaml.Unmarshal(readFileOrPanic(spec))
		if err != nil {
			panic(catch.Error(err))
		}
		vMap, ok := m.(dgo.Map)
		if !ok {
			panic(catch.Error(`expecting data to be a map`))
		}
		sType = tf.StructMapFromMap(false, vMap)
	case strings.HasSuffix(spec, `.dgo`):
		tp := tf.ParseFile(nil, spec, string(readFileOrPanic(spec)))
		if st, ok := tp.(dgo.StructMapType); ok {
			sType = st
		} else {
			panic(catch.Error(`file '%s' does not contain a struct definition`, spec))
		}
	default:
		panic(catch.Error(`invalid file name '%s', expected file name to end with .yaml, .json, or .dgo`, spec))
	}
	return
}

// Do parses the validate command line options and runs the validation
func (h *validateCommand) Do(args []string) int {
	return h.RunWithCatch(func() int {
		r, done := h.Parse(args)
		if done {
			return r
		}
		if h.input == `` {
			return h.MissingOption(`input`)
		}
		if h.spec == `` {
			return h.MissingOption(`spec`)
		}
		return h.run()
	})
}
