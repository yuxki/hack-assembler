package hack

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	noSymbolAddCommands = `
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/06/add/Add.asm

// Computes R0 = 2 + 3  (R0 refers to RAM[0])

@2
D=A
@3
D=D+A
@0
M=D
`

	noSymbolAddCommandsBinary = `0000000000000010
1110110000010000
0000000000000011
1110000010010000
0000000000000000
1110001100001000
`
)

const (
	noSymbolMaxCommands = `
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/06/max/MaxL.asm

// Symbol-less version of the Max.asm program.

@0
D=M
@1
D=D-M
@12
D;JGT
@1
D=M
@2
M=D
@16
0;JMP
@0
D=M
@2
M=D
@16
0;JMP
`

	noSymbolMaxCommandsBinary = `0000000000000000
1111110000010000
0000000000000001
1111010011010000
0000000000001100
1110001100000001
0000000000000001
1111110000010000
0000000000000010
1110001100001000
0000000000010000
1110101010000111
0000000000000000
1111110000010000
0000000000000010
1110001100001000
0000000000010000
1110101010000111
`
)

const (
	noSymbolRectCommands = `
// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/06/rect/RectL.asm

// Symbol-less version of the Rect.asm program.

@0
D=M
@24
D;JLE
@16
M=D
@16384
D=A
@17
M=D
@17
A=M
M=-1
@17
D=M
@32
D=D+A
@17
M=D
@16
M=M-1
D=M
@10
D;JGT
@24
0;JMP
`

	noSymbolRectCommandsBinary = `0000000000000000
1111110000010000
0000000000011000
1110001100000110
0000000000010000
1110001100001000
0100000000000000
1110110000010000
0000000000010001
1110001100001000
0000000000010001
1111110000100000
1110111010001000
0000000000010001
1111110000010000
0000000000100000
1110000010010000
0000000000010001
1110001100001000
0000000000010000
1111110010001000
1111110000010000
0000000000001010
1110001100000001
0000000000011000
1110101010000111
`
)

func TestAssembler_Assemble(t *testing.T) {
	data := []struct {
		testCase string
		asm      string
		binary   string
	}{
		{
			testCase: "add",
			asm:      noSymbolAddCommands,
			binary:   noSymbolAddCommandsBinary,
		},
		{
			testCase: "max",
			asm:      noSymbolMaxCommands,
			binary:   noSymbolMaxCommandsBinary,
		},
		{
			testCase: "rect",
			asm:      noSymbolRectCommands,
			binary:   noSymbolRectCommandsBinary,
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.testCase, func(t *testing.T) {
			parser := NewParser(strings.NewReader(d.asm))
			code := NewCode()
			writer := &bytes.Buffer{}

			assembler := NewAssembler(writer, parser, code)
			err := assembler.Assemble()
			if err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(writer.String(), d.binary); diff != "" {
				print(writer.String())
				t.Errorf(diff)
			}
		})
	}
}