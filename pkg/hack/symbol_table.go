package hack

import (
	"errors"
	"fmt"
	"regexp"
)

// Entry represents a symbol table entry.
type Entry struct {
	symbol  string
	address uint
}

const (
	screenAddress uint = 16384
)

// ErrInvalidSymbol is returned when the symbol is invalid.
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

// SymbolTable is a data structure that is used to map symbolic labels or
// variables to their corresponding numeric addresses. It serves as a lookup table for
// retrieving addresses based on symbol names.
type SymbolTable struct {
	entries []Entry
}

const (
	initialTableCapacity = 23
)

// NewSymbolTable creates a new symbol table.
// It initializes the table with the predefined symbols.
func NewSymbolTable() *SymbolTable {
	table := &SymbolTable{make([]Entry, 0, initialTableCapacity)}
	return table
}

var ErrSymbolAlreadyExists = errors.New("symbol already exists")

// AddEntry adds the pair (symbol, address) to the table.
// It returns an error if the symbol already exists or if the symbol is invalid.
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

// Contains returns true if the symbol table contains the given symbol.
func (s *SymbolTable) Contains(symbol string) bool {
	for _, entry := range s.entries {
		if entry.symbol == symbol {
			return true
		}
	}

	return false
}

var ErrSymbolNotFound = errors.New("symbol not found")

// GetAddress returns the address associated with the symbol.
// It returns an error if the symbol is not found.
func (s *SymbolTable) GetAddress(symbol string) (uint, error) {
	for _, entry := range s.entries {
		if entry.symbol == symbol {
			return entry.address, nil
		}
	}

	return 0, fmt.Errorf("could not get address: %w", ErrSymbolNotFound)
}
