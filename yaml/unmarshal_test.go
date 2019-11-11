package yaml_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lyraproj/dgo/dgo"

	require "github.com/lyraproj/dgo/dgo_test"
	"github.com/lyraproj/dgo/typ"
	"github.com/lyraproj/dgo/vf"
	"github.com/lyraproj/dgoyaml/yaml"
)

func ExampleUnmarshal() {
	v, err := yaml.Unmarshal([]byte(`
- hello
- true
- 1
- 3.14
- null
- a: 1`))
	if err == nil {
		fmt.Println(v)
	}
	// Output: {"hello",true,1,3.14,nil,{"a":1}}
}

func TestUnmarshal_map(t *testing.T) {
	m, err := yaml.Unmarshal([]byte(`
a: 1
b: two
c: 
  - hello
  - true
  - 1
  - 3.14
  - null
`))
	require.Ok(t, err)
	require.Equal(t, vf.Map("a", 1, "b", "two", "c", vf.Values(`hello`, true, 1, 3.14, nil)), m)
}

func TestUnmarshal_binary(t *testing.T) {
	m, err := yaml.Unmarshal([]byte("b: !!binary AQQD\n"))
	require.Ok(t, err)
	require.Equal(t, vf.Map("b", vf.BinaryFromString(`AQQD`)), m)
	require.Panic(t, func() { _, _ = yaml.Unmarshal([]byte("b: !!binary AQQ~\n")) }, `illegal base64 data`)
}

func TestUnmarshal_timestamp(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, `2019-10-06T07:15:00-07:00`)
	m, err := yaml.Unmarshal([]byte("t: !!timestamp 2019-10-06T07:15:00-07:00\n"))
	require.Ok(t, err)
	require.Equal(t, vf.Map("t", vf.Time(ts)), m)

	_, err = yaml.Unmarshal([]byte("t: !!timestamp 2019-13-06T07:15:00-07:00\n"))
	require.NotOk(t, `cannot decode`, err)
}

func TestUnmarshal_bad_yaml(t *testing.T) {
	_, err := yaml.Unmarshal([]byte(": :\n"))
	require.NotOk(t, `did not find expected key`, err)
}

func TestUnmarshal_type(t *testing.T) {
	m, err := yaml.Unmarshal([]byte("t: !puppet.com,2019:dgo/type string\n"))
	require.Ok(t, err)
	require.Equal(t, vf.Map("t", typ.String), m)
}

func TestUnmarshal_typedMap(t *testing.T) {
	_, err := yaml.Unmarshal([]byte(`
int: 23
float: 3.14`))
	require.Ok(t, err)

	require.Panic(t, func() {
		m, err := yaml.Unmarshal([]byte(`
int: 23
string: hello`))
		require.Ok(t, err)
		m.(dgo.Map).SetType(`map[string](int|float)`)
	}, `the string "hello" cannot be assigned to a variable of type int|float`)
}

func TestUnmarshal_array(t *testing.T) {
	a, err := yaml.Unmarshal([]byte(`["hello",true,1,3.14,null]`))
	require.Ok(t, err)
	require.Equal(t, vf.Values(`hello`, true, 1, 3.14, nil), a)
}
