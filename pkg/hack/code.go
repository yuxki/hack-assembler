package hack

import (
	"errors"
	"fmt"
)

// Code is a struct that translates Hack assembly language mnemonics into binary codes.
type Code struct{}

// NewCode returns a new Code.
func NewCode() Code {
	return Code{}
}

// ErrInvalidNemonic is returned when the nemonic is invalid.
var ErrInvalidNemonic = errors.New("invalid nemonic")

// Dest returns the binary code of the dest mnemonic.
func (c Code) Dest(n string) (string, error) {
	if n == "" {
		return "000", nil
	}
	if n == "M" {
		return "001", nil
	}
	if n == "D" {
		return "010", nil
	}
	if n == "MD" {
		return "011", nil
	}
	if n == "A" {
		return "100", nil
	}
	if n == "AM" {
		return "101", nil
	}
	if n == "AD" {
		return "110", nil
	}
	if n == "AMD" {
		return "111", nil
	}

	return "", fmt.Errorf("could not convert a dest command:%s: %w", n, ErrInvalidNemonic)
}

// Comp returns the binary code of the comp mnemonic.
func (c Code) Comp(n string) (string, error) {
	if n == "0" {
		return "0101010", nil
	}
	if n == "1" {
		return "0111111", nil
	}
	if n == "-1" {
		return "0111010", nil
	}
	if n == "D" {
		return "0001100", nil
	}
	if n == "A" {
		return "0110000", nil
	}
	if n == "!D" {
		return "0001101", nil
	}
	if n == "!A" {
		return "0110001", nil
	}
	if n == "-D" {
		return "0001111", nil
	}
	if n == "-A" {
		return "0110011", nil
	}
	if n == "D+1" {
		return "0011111", nil
	}
	if n == "A+1" {
		return "0110111", nil
	}
	if n == "D-1" {
		return "0001110", nil
	}
	if n == "A-1" {
		return "0110010", nil
	}
	if n == "D+A" {
		return "0000010", nil
	}
	if n == "D-A" {
		return "0010011", nil
	}
	if n == "A-D" {
		return "0000111", nil
	}
	if n == "D&A" {
		return "0000000", nil
	}
	if n == "D|A" {
		return "0010101", nil
	}
	if n == "M" {
		return "1110000", nil
	}
	if n == "!M" {
		return "1110001", nil
	}
	if n == "-M" {
		return "1110011", nil
	}
	if n == "M+1" {
		return "1110111", nil
	}
	if n == "M-1" {
		return "1110010", nil
	}
	if n == "D+M" {
		return "1000010", nil
	}
	if n == "D-M" {
		return "1010011", nil
	}
	if n == "M-D" {
		return "1000111", nil
	}
	if n == "D&M" {
		return "1000000", nil
	}
	if n == "D|M" {
		return "1010101", nil
	}

	return "", fmt.Errorf("could not convert a comp command:%s: %w", n, ErrInvalidNemonic)
}

// Jump returns the binary code of the jump mnemonic.
func (c Code) Jump(n string) (string, error) {
	if n == "" {
		return "000", nil
	}
	if n == "JGT" {
		return "001", nil
	}
	if n == "JEQ" {
		return "010", nil
	}
	if n == "JGE" {
		return "011", nil
	}
	if n == "JLT" {
		return "100", nil
	}
	if n == "JNE" {
		return "101", nil
	}
	if n == "JLE" {
		return "110", nil
	}
	if n == "JMP" {
		return "111", nil
	}

	return "", fmt.Errorf("could not convert a jump command:%s: %w", n, ErrInvalidNemonic)
}
