package hack

import (
	"errors"
	"fmt"
	"regexp"
)

type Entry struct {
	symbol  string
	address uint
}

const (
	screenAddress uint = 16384
)

var ErrInvalidSymbol = errors.New("invalid symbol")

func newEntry(symbol string, address uint) (Entry, error) {
	var entry Entry
	ok, err := regexp.MatchString("^[a-zA-Z_.$:][a-zA-Z0-9_.$:]*$", symbol)
	if err != nil {
		return entry, err
	}

	if !ok {
		return entry, fmt.Errorf("could not create entry: %w", ErrInvalidSymbol)
	}

	entry.symbol = symbol
	entry.address = address

	return entry, nil
}

type SymbolTable struct {
	entries []Entry
}

func NewSymbolTable() (*SymbolTable, error) {
	table := &SymbolTable{make([]Entry, 0, 23)}

	err := table.AddEntry("SP", 0)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("LCL", 1)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("ARG", 2)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("THIS", 3)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("THAT", 4)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("SCREEN", screenAddress)
	if err != nil {
		return nil, err
	}

	err = table.AddEntry("KBD", 24576)
	if err != nil {
		return nil, err
	}

	var i uint
	for i = 0; i < 16; i++ {
		err = table.AddEntry(fmt.Sprintf("R%d", i), i)
		if err != nil {
			return nil, err
		}
	}

	return table, nil
}

var ErrSymbolAlreadyExists = errors.New("symbol already exists")

func (s *SymbolTable) AddEntry(symbol string, address uint) error {
	if s.Contains(symbol) {
		return fmt.Errorf("could not add entry to the sybmol table: %w", ErrSymbolAlreadyExists)
	}

	entry, err := newEntry(symbol, address)
	if err != nil {
		return fmt.Errorf("could not add entry to the sybmol table: %w", err)
	}

	s.entries = append(s.entries, entry)

	return nil
}

func (s *SymbolTable) Contains(symbol string) bool {
	for _, entry := range s.entries {
		if entry.symbol == symbol {
			return true
		}
	}

	return false
}

var ErrSymbolNotFound = errors.New("symbol not found")

func (s *SymbolTable) GetAddress(symbol string) (uint, error) {
	for _, entry := range s.entries {
		if entry.symbol == symbol {
			return entry.address, nil
		}
	}

	return 0, fmt.Errorf("could not get address: %w", ErrSymbolNotFound)
}
