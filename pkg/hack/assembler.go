package hack

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

type Assembler struct {
	w      io.Writer
	parser *Parser
	code   Code
}

func NewAssembler(w io.Writer, parser *Parser, code Code) *Assembler {
	return &Assembler{w: w, parser: parser, code: code}
}

var ErrInvalidCommand = errors.New("invalid command")

func (a *Assembler) Assemble() error {
	for a.parser.Advance() {
		binary := ""
		switch a.parser.CommandType() {
		case ACommand:
			symbol, err := a.parser.Symbol()
			if err != nil {
				return err
			}
			if regexp.MustCompile(`^\d+$`).MatchString(symbol) {
				symInt, err := strconv.Atoi(symbol)
				if err != nil {
					return err
				}
				binary = fmt.Sprintf("0%015b", symInt)
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
