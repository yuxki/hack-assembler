package hack

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Assembler is a struct that assembles Hack assembly code into Hack machine code.
// It uses a Parser to parse the assembly code, a Code to translate the parsed commands into binary,
// and a SymbolTable to keep track of symbols and their addresses.
type Assembler struct {
	w           io.Writer
	parser      *Parser
	code        Code
	symbolTable *SymbolTable
	nextAddress uint
}

const (
	spAddress   = 0
	lclAddress  = 1
	argAddress  = 2
	thisAddress = 3
	thatAddress = 4
	kbdAddress  = 24576

	initialNextAddress = 16

	ramAddressNumbers = initialNextAddress
)

// NewAssembler creates a new instance of the Assembler.
// It takes a reader and a writer as input parameters.
// The reader is used to read the assembly code, while the writer is used to write the machine code.
// It returns a pointer to the Assembler and an error (if any) occurred during initialization.
func NewAssembler(r io.Reader, w io.Writer) (*Assembler, error) {
	parser := NewParser(r)
	code := NewCode()
	table := NewSymbolTable()

	err := table.AddEntry("SP", spAddress)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("LCL", lclAddress)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("ARG", argAddress)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("THIS", thisAddress)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("THAT", thatAddress)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("SCREEN", screenAddress)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("KBD", kbdAddress)
	if err != nil {
		return nil, err
	}

	var i uint
	for i = 0; i < ramAddressNumbers; i++ {
		err = table.AddEntry(fmt.Sprintf("R%d", i), i)
		if err != nil {
			return nil, err
		}
	}

	return &Assembler{
			w:           w,
			parser:      parser,
			code:        code,
			symbolTable: table,
			nextAddress: initialNextAddress,
		},
		nil
}

// ErrInvalidCommand is returned when the parser encounters an invalid command.
var ErrInvalidCommand = errors.New("invalid command")

func uintAddressToACommandBinary(address uint) string {
	return fmt.Sprintf("0%015b", address)
}

func intAddressToACommandBinary(address int) string {
	return fmt.Sprintf("0%015b", address)
}

func (a *Assembler) assembleACommand() (string, error) {
	var binary string

	symbol, err := a.parser.Symbol()
	if err != nil {
		return "", err
	}
	if regexp.MustCompile(`^\d+$`).MatchString(symbol) {
		address, err := strconv.Atoi(symbol)
		if err != nil {
			return "", err
		}
		return intAddressToACommandBinary(address), nil
	}

	if a.symbolTable.Contains(symbol) {
		address, err := a.symbolTable.GetAddress(symbol)
		if err != nil {
			return "", err
		}
		return uintAddressToACommandBinary(address), nil
	}

	err = a.symbolTable.AddEntry(symbol, a.nextAddress)
	if err != nil {
		return "", err
	}
	binary = uintAddressToACommandBinary(a.nextAddress)
	a.nextAddress++

	return binary, nil
}

func (a *Assembler) assembleCCommand() (string, error) {
	dest, err := a.parser.Dest()
	if err != nil {
		return "", err
	}
	dBits, err := a.code.Dest(dest)
	if err != nil {
		return "", err
	}

	comp, err := a.parser.Comp()
	if err != nil {
		return "", err
	}
	cBits, err := a.code.Comp(comp)
	if err != nil {
		return "", err
	}

	jump, err := a.parser.Jump()
	if err != nil {
		return "", err
	}
	jBits, err := a.code.Jump(jump)
	if err != nil {
		return "", err
	}

	if dest+jump == "" {
		return "", fmt.Errorf("no dest or jump: %w", ErrInvalidCommand)
	}

	return fmt.Sprintf("111%s%s%s", cBits, dBits, jBits), nil
}

// Assemble function takes the Hack assembly code as input and converts it
// into Hack machine code.
// It then writes the machine code to the writer provided by the Assembler.
// If the assembly code is invalid, it will return an error.
// To accomplish this, it performs the following steps:
// 1. Creation of a symbol table.
// 2. Parsing of the assembly code.
// 3. Translation of the parsed code into binary.
func (a *Assembler) Assemble() error {
	err := a.createSymbolTable()
	if err != nil {
		return err
	}

	for a.parser.Advance() {
		binary := ""
		switch a.parser.CommandType() {
		case ACommand:
			binary, err = a.assembleACommand()
			if err != nil {
				return err
			}
		case CCommand:
			binary, err = a.assembleCCommand()
			if err != nil {
				return err
			}
		case LCommand:
			continue
		}
		if binary == "" {
			return fmt.Errorf("could not assemble binary with %s: %w", a.parser.Command(), ErrInvalidCommand)
		}
		_, err := a.w.Write([]byte(binary + "\n"))
		if err != nil {
			return err
		}
	}
	return nil
}

// createSymbolTable function creates a symbol table from the assembly code.
// It does this by parsing the assembly code and adding symbols to the symbol table
// as they are encountered.
func (a *Assembler) createSymbolTable() error {
	commands := ""

	for a.parser.Advance() {
		commands += a.parser.Command() + "\n"

		switch a.parser.CommandType() {
		case LCommand:
			symbol, err := a.parser.Symbol()
			if err != nil {
				return err
			}
			err = a.symbolTable.AddEntry(symbol, a.parser.LineNumber())
			if err != nil {
				return err
			}
		case ACommand:
		case CCommand:
			continue
		}
	}

	a.parser.Reset(strings.NewReader(commands))

	return nil
}
