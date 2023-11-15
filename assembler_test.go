package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	noSymbolCommands = `@2
D=A
@3
D=D+A
@0
M=D`

	noSymbolCommandsBinary = `0000000000000010
1110110000010000
0000000000000011
1110000010010000
0000000000000000
1110001100001000
`
)

func TestAssembler_Assemble(t *testing.T) {
	t.Parallel()

	parser := NewParser(strings.NewReader(noSymbolCommands))
	code := NewCode()
	writer := &bytes.Buffer{}

	assembler := NewAssembler(writer, parser, code)
	err := assembler.Assemble()
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(writer.String(), noSymbolCommandsBinary); diff != "" {
		t.Errorf(diff)
	}
}
