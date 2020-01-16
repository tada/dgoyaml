package yaml_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"
	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
	"github.com/lyraproj/dgoyaml/yaml"
	y3 "gopkg.in/yaml.v3"
)

func TestMarshal(t *testing.T) {
	m := vf.Map("a", 1, "b", "two", "c", vf.Values(`hello`, true, 1, 3.14, nil))
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `a: 1
b: two
c:
  - hello
  - true
  - 1
  - 3.14
  - null
`, string(b))
}

func TestMarshal_binary(t *testing.T) {
	m := vf.Map("b", vf.BinaryFromString(`AQQD`))
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `b: !!binary AQQD
`, string(b))
}

func TestMarshal_structMap(t *testing.T) {
	type structA struct {
		A string `json:"a"`
		B int    `json:"b"`
	}

	type structAyaml struct {
		A string `yaml:"a"`
		B int    `yaml:"b"`
	}

	s := structA{A: `Alpha`, B: 32}
	m := vf.Map(&s)
	j, err := yaml.Marshal(m)
	require.Ok(t, err)
	require.Equal(t, `a: Alpha
b: 32
`, string(j))

	s2 := structAyaml{A: `Alpha`, B: 32}
	m = vf.Map(&s2)
	j, err = yaml.Marshal(m)
	require.Ok(t, err)
	require.Equal(t, `a: Alpha
b: 32
`, string(j))

	m = vf.Map(`nested`, m)
	j, err = yaml.Marshal(m)
	require.Ok(t, err)
	require.Equal(t, `nested:
    a: Alpha
    b: 32
`, string(j))

	type structFailYaml struct {
		A string           `yaml:"a"`
		B *marshalTestFail `yaml:"b"`
	}
	sFail := structFailYaml{A: `Alpha`, B: &marshalTestFail{}}
	m = vf.Map(&sFail)
	_, err = yaml.Marshal(m)
	require.NotNil(t, err)
	require.Equal(t, `errFailing`, err.Error())
}

func TestMarshal_timestamp(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	m := vf.Map("t", vf.Time(ts))
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `t: !!timestamp 2019-10-06T07:15:00-07:00
`, string(b))
}

func TestMarshal_type(t *testing.T) {
	m := vf.Map("t", typ.String)
	b, err := yaml.Marshal(m)
	require.Nil(t, err)
	require.Equal(t, `t: !puppet.com,2019:dgo/type string
`, string(b))
}

var errFailing = errors.New("errFailing")

type testNoMarshaler struct {
	A string
}

type testMarshaler struct {
	A string
}

type marshalTestNode struct {
}

type marshalTestFail struct {
}

type marshalTestPanic struct {
}

func (m *testMarshaler) MarshalYAML() (interface{}, error) {
	return vf.Map(`A`, m.A), nil
}

func (m *marshalTestFail) MarshalYAML() (interface{}, error) {
	return nil, errFailing
}

func (m *marshalTestPanic) MarshalYAML() (interface{}, error) {
	panic(errFailing)
}

func (m *marshalTestNode) MarshalYAML() (interface{}, error) {
	return &y3.Node{Kind: y3.ScalarNode, Tag: `!!int`, Value: `23`}, nil
}

// obscureValue is to get test coverage of unknown value types
type obscureValue int

func (o obscureValue) String() string {
	panic("implement me")
}

func (o obscureValue) Type() dgo.Type {
	return typ.Any
}

func (o obscureValue) Equals(other interface{}) bool {
	panic("implement me")
}

func (o obscureValue) HashCode() int {
	panic("implement me")
}

func TestMarshal_native(t *testing.T) {
	m := vf.MutableValues(&testMarshaler{A: `hello`})
	bs, err := yaml.Marshal(m)
	require.Ok(t, err)
	require.Equal(t, "- A: hello\n", string(bs))
}

func TestMarshal_native_node(t *testing.T) {
	m := vf.MutableValues(&marshalTestNode{})
	bs, err := yaml.Marshal(m)
	require.Ok(t, err)
	require.Equal(t, "- 23\n", string(bs))
}

func TestMarshal_obscure(t *testing.T) {
	m := vf.MutableValues(obscureValue(0))
	_, err := yaml.Marshal(m)
	require.NotNil(t, err)
	require.Equal(t, `unable to marshal into value of type any`, err.Error())
}

func TestMarshal_fail(t *testing.T) {
	m := vf.MutableValues(&marshalTestFail{})
	_, err := yaml.Marshal(m)
	require.NotNil(t, err)
	require.Equal(t, `errFailing`, err.Error())
}

func TestMarshal_panic(t *testing.T) {
	m := vf.MutableValues(&marshalTestPanic{})
	require.Panic(t, func() { _, _ = yaml.Marshal(m) }, `errFailing`)
}

func TestMarshal_failNoMarshaler(t *testing.T) {
	m := vf.MutableValues(&testNoMarshaler{})
	_, err := yaml.Marshal(m)
	require.NotNil(t, err)
	require.Equal(t, `unable to marshal into value of type *yaml_test.testNoMarshaler`, err.Error())
}
