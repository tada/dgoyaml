# dgoyaml Dgo YAML bi-directional serialization

[![](https://goreportcard.com/badge/github.com/lyraproj/dgoyaml)](https://goreportcard.com/report/github.com/lyraproj/dgoyaml)
[![](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/lyraproj/dgoyaml)
[![](https://github.com/lyraproj/dgoyaml/workflows/Dgo%20YAML%20Build/badge.svg)](https://github.com/lyraproj/dgoyaml/actions)

This module provides YAML serialization/deserialization for [Dgo](https://github.com/lyraproj/dgo) using the
[gopkg.in/yaml.v3](https://github.com/go-yaml/yaml/tree/v3) module.

### Using dgoyaml as a library
To use dgoyaml, first install the latest version of the library:
```sh
go get github.com/lyraproj/dgoyaml
```

### Running the dgo CLI
The dgo CLI command can be used to get acquainted with dgo concepts. It will allow you to declare
types and values in YAML files and then use the types to validate the values. You install the
command under $GOPATH/bin with:
```sh
go install github.com/lyraproj/dgoyaml/cli/dgo
```
after that, you should be able to do:
```sh
dgo help
```
to get a description of avaliable sub commands and flags.

### Example of usage:
Let's assume some kind of typed parameters in YAML that the user enters like this:
```yaml
host: example.com
port: 22
```
The task is to create a user friendly description of a parameter, also in YAML, which can be used to validate
the above parameters. Something like this:
```yaml
host:
  type: string[1]
  name: sample/service_host
  required: true
port:
  type: 1..999
  name: sample/service_port
```
The value of each `type` is a [dgo type](docs/types.md). They limit the host parameter to a non empty string
and the port parameter to an integer in the range 1-999. A special `required` entry is used to denote whether
or not a parameter value must be present. The `name` entry is optional and provides a freeform text identifier.

Put the two above YAML examples in two separate files, `params.yaml` and `params_spec.yaml`. Then run the
command:
```
dgo validate --verbose --input params.yaml --spec params_spec.yaml
```
The output should be:
```
Got input yaml with:
  host: example.com
  port: 22
Validating 'host' against definition string[1]
  'host' OK!
Validating 'port' against definition 1..999
  'port' OK!
```
For examples of how to use the library functions that Dgo provides to perform the above validation in, please take
a look at [parameter_test.go](examples_test/parameter_test.go). The source of the [validate command](cli/validate.go)
may also be of help.
