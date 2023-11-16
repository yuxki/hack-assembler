package hack

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type Assembler struct {
	w           io.Writer
	parser      *Parser
	code        Code
	symbolTable *SymbolTable
	nextAddress uint
}

func NewAssembler(w io.Writer, parser *Parser, code Code, symbolTable *SymbolTable) *Assembler {
	return &Assembler{w: w, parser: parser, code: code, symbolTable: symbolTable, nextAddress: 16}
}

var ErrInvalidCommand = errors.New("invalid command")

func (a *Assembler) Assemble() error {
	err := a.createSymbolTable()
	if err != nil {
		return err
	}

	for a.parser.Advance() {
		binary := ""
		switch a.parser.CommandType() {
		case ACommand:
			symbol, err := a.parser.Symbol()
			if err != nil {
				return err
			}
			if regexp.MustCompile(`^\d+$`).MatchString(symbol) {
				address, err := strconv.Atoi(symbol)
				if err != nil {
					return err
				}
				binary = fmt.Sprintf("0%015b", address)
			} else {
				if a.symbolTable.Contains(symbol) {
					address, err := a.symbolTable.GetAddress(symbol)
					if err != nil {
						return err
					}
					binary = fmt.Sprintf("0%015b", address)
				} else {
					err := a.symbolTable.AddEntry(symbol, a.nextAddress)
					if err != nil {
						return err
					}
					binary = fmt.Sprintf("0%015b", a.nextAddress)
					a.nextAddress++
				}
			}
		case CCommand:
			dest, err := a.parser.Dest()
			if err != nil {
				return err
			}
			dBits, err := a.code.Dest(dest)
			if err != nil {
				return err
			}

			comp, err := a.parser.Comp()
			if err != nil {
				return err
			}
			cBits, err := a.code.Comp(comp)
			if err != nil {
				return err
			}

			jump, err := a.parser.Jump()
			if err != nil {
				return err
			}
			jBits, err := a.code.Jump(jump)
			if err != nil {
				return err
			}

			if dest+jump == "" {
				return fmt.Errorf("no dest or jump: %w", ErrInvalidCommand)
			}

			binary = fmt.Sprintf("111%s%s%s", cBits, dBits, jBits)

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
		default:
			continue
		}
	}

	a.parser.Reset(strings.NewReader(commands))

	return nil
}
