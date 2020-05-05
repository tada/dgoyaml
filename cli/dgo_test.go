package cli_test

import (
	"strings"
	"testing"

	"github.com/tada/dgo/test/assert"
	"github.com/tada/dgoyaml/cli"
)

func TestDgo_noArgs(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{}))
	assert.Match(t, `missing required command`, err.String())
}

func TestDgo_unknownSubCommand(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`what`}))
	assert.Match(t, `unknown command: what`, err.String())
}

func TestDgo_unknownFlag(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`-what`}))
	assert.Match(t, `flag provided but not defined: -what`, err.String())
}

func TestDgo_validate_missingOption(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`validate`}))
	assert.Match(t, `missing required option: -input`, err.String())
	assert.Equal(t, 1, dgo.Do([]string{`validate`, `--input`, `foo`}))
	assert.Match(t, `missing required option: -spec`, err.String())
	assert.Equal(t, 1, dgo.Do([]string{`validate`, `--spec`, `foo`}))
	assert.Match(t, `missing required option: -input`, err.String())
}

func TestDgo_validate_missingOptionArgument(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`validate`, `--input`}))
	assert.Match(t, `flag needs an argument: -input`, err.String())
	assert.Equal(t, 1, dgo.Do([]string{`validate`, `--input`, `foo`, `--spec`}))
	assert.Match(t, `flag needs an argument: -spec`, err.String())
}

func TestDgo_validate_unknownOption(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`validate`, `--output`}))
	assert.Match(t, `flag provided but not defined: -output`, err.String())
}

func TestDgo_help(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 0, dgo.Do([]string{`help`}))
	assert.Match(t, `dgo: a command `, out.String())
	assert.Equal(t, 0, dgo.Do([]string{`--help`}))
	assert.Match(t, `dgo: a command `, out.String())
}

func TestDgo_validate_noArgs(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`validate`}))
	assert.Match(t, `missing required option: -input`, err.String())
}

func TestDgo_validate_help(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 0, dgo.Do([]string{`help`, `validate`}))
	assert.Match(t, `dgo validate\s+-input `, out.String())
}

func TestDgo_validate_ok(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 0, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	s := out.String()
	assert.Match(t, `'host' OK\!`, s)
	assert.Match(t, `'port' OK\!`, s)
}

func TestDgo_validate_ok_brief(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 0, dgo.Do([]string{`validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	assert.Equal(t, ``, out.String())
}

func TestDgo_validate_bad_port(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_bad_port.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	s := out.String()
	assert.NoMatch(t, `'host' FAILED\!`, s)
	assert.Match(t, `'port' FAILED\!(?:.|\s)*2222`, s)
}

func TestDgo_validate_bad_port_brief(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`validate`, `--input`, `testdata/service_bad_port.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	assert.Equal(t, "parameter 'port' is not an instance of type 1..999\n", out.String())
}

func TestDgo_validate_bad_port_dgo(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_bad_port.yaml`, `--spec`, `testdata/servicespec.dgo`}))
	s := out.String()
	assert.NoMatch(t, `'host' FAILED\!`, s)
	assert.Match(t, `'port' FAILED\!(?:.|\s)*2222`, s)
}

func TestDgo_validate_extraneous_param(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_extraneous_param.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	s := out.String()
	assert.Match(t, `'host' OK\!`, s)
	assert.Match(t, `'port' OK\!`, s)
	assert.Match(t, `'login' FAILED\!(?:.|\s)*key is not found in definition`, s)
}

func TestDgo_validate_missing_host(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_missing_host.yaml`, `--spec`, `testdata/servicespec.yaml`}))
	s := out.String()
	assert.Match(t, `'port' OK\!`, s)
	assert.Match(t, `'host' FAILED\!(?:.|\s)*required key not found in input`, s)
}

func TestDgo_validate_no_input_file(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/server.yaml`, `--spec`, `testdata/servicespec.dgo`}))
	s := err.String()
	assert.Match(t, `server\.yaml.*no such file or directory`, s)
}

func TestDgo_validate_no_spec_file(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/serverspec.dgo`}))
	s := err.String()
	assert.Match(t, `serverspec\.dgo.*no such file or directory`, s)
}

func TestDgo_validate_input_not_map(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service_array.yaml`, `--spec`, `testdata/servicespec.dgo`}))
	s := err.String()
	assert.Match(t, `Error: expecting data to be a map`, s)
}

func TestDgo_validate_spec_not_map(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec_array.yaml`}))
	s := err.String()
	assert.Match(t, `Error: expecting data to be a map`, s)
}

func TestDgo_validate_spec_bad_yaml(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/bad.yaml`}))
	s := err.String()
	assert.Match(t, `did not find expected key`, s)
}

func TestDgo_validate_spec_bad_type(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec_bad_type.dgo`}))
	s := err.String()
	assert.Match(t, `does not contain a struct definition`, s)
}

func TestDgo_validate_input_extension(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.pson`, `--spec`, `testdata/servicespec.dgo`}))
	s := err.String()
	assert.Match(t, `expected file name to end with \.yaml or \.json`, s)
}

func TestDgo_validate_spec_extension(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec.go`}))
	s := err.String()
	assert.Match(t, `expected file name to end with \.yaml, \.json, or \.dgo`, s)
}

func TestDgo_validate_spec_bad_dgo(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/service.yaml`, `--spec`, `testdata/servicespec_bad.dgo`}))
	s := err.String()
	assert.Match(t, `Error: mix of elements and map entries`, s)
}

func TestDgo_validate_bad_yaml(t *testing.T) {
	out := &strings.Builder{}
	err := &strings.Builder{}
	dgo := cli.Dgo(out, err)
	assert.Equal(t, 1, dgo.Do([]string{`--verbose`, `validate`, `--input`, `testdata/bad.yaml`, `--spec`, `testdata/servicespec.go`}))
	s := err.String()
	assert.Match(t, `did not find expected key`, s)
}
