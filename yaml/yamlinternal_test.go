package yaml

import (
	"testing"

	require "github.com/lyraproj/dgo/dgo_test"
	"gopkg.in/yaml.v3"
)

func badNode() *yaml.Node {
	n := &yaml.Node{Kind: yaml.AliasNode}
	n.Alias = n
	return n
}

func unknownTagNode() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: `!something:here`, Value: `a`}
}

func Test_decodeScalar_unknownTagNode(t *testing.T) {
	require.Equal(t, `a`, decodeScalar(unknownTagNode()))
}

func Test_decodeScalar_badNode(t *testing.T) {
	require.Panic(t, func() { decodeScalar(badNode()) }, `value contains itself`)
}
