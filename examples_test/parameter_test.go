package examples

import (
	"testing"

	"github.com/tada/dgo/dgo"
	"github.com/tada/dgo/tf"
	"github.com/tada/dgoyaml/yaml"
)

// Sample parameter map
const sampleParameters = `
host:
  type: string[1]
  name: sample/service_host
  required: true
port:
  type: 1..999
  name: sample/service_port
`

func TestValidateParameterValues(t *testing.T) {
	const sampleValues = `
host: example.com
port: 22
`
	params, err := yaml.Unmarshal([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectNoErrors(t, validate(t, params))
}

func TestValidateParameterValues_failRequired(t *testing.T) {
	const sampleValues = `
port: 22
`
	params, err := yaml.Unmarshal([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectError(t, `missing required parameter 'host'`, validate(t, params))
}

func TestValidateParameterValues_failNotRecognized(t *testing.T) {
	const sampleValues = `
host: example.com
port: 22
login: foo:bar
`
	params, err := yaml.Unmarshal([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectError(t, `unknown parameter 'login'`, validate(t, params))
}

func TestValidateParameterValues_failInvalidHostType(t *testing.T) {
	const sampleValues = `
host: 85493
port: 22
`
	params, err := yaml.Unmarshal([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectError(t, `parameter 'host' is not an instance of type string[1]`, validate(t, params))
}

func TestValidateParameterValues_failInvalidPortType(t *testing.T) {
	const sampleValues = `
host: example.com
port: 1022
`
	params, err := yaml.Unmarshal([]byte(sampleValues))
	if err != nil {
		t.Fatal(err)
	}
	expectError(t, `parameter 'port' is not an instance of type 1..999`, validate(t, params))
}

func validate(t *testing.T, params dgo.Value) []error {
	t.Helper()
	pt, err := loadDesc([]byte(sampleParameters))
	if err != nil {
		t.Fatal(err)
	}
	return pt.(dgo.MapValidation).Validate(nil, params)
}

func loadDesc(yamlData []byte) (dgo.StructMapType, error) {
	data, err := yaml.Unmarshal(yamlData)
	if err != nil {
		return nil, err
	}
	return tf.StructMapFromMap(false, data.(dgo.Map)), nil
}

func expectError(t *testing.T, error string, errors []error) {
	t.Helper()
	if len(errors) != 1 || error != errors[0].Error() {
		t.Errorf(`expected "%s" error`, error)
		expectNoErrors(t, errors)
	}
}

func expectNoErrors(t *testing.T, errors []error) {
	t.Helper()
	for _, err := range errors {
		t.Error(err)
	}
}
