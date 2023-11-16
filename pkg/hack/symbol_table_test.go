package hack

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSymbolTable_AddEntry_Contains_GetAddress(t *testing.T) {
	t.Parallel()

	table := NewSymbolTable()

	data := []struct {
		testCase string
		symbol   string
		address  uint
		err      error
	}{
		{
			testCase: "AddEntry: valid symbol",
			symbol:   "test",
			address:  0,
			err:      nil,
		},
		{
			testCase: "AddEntry: invalid symbol",
			symbol:   "2",
			address:  0,
			err:      ErrInvalidSymbol,
		},
		{
			testCase: "AddEntry: symbol already exists",
			symbol:   "test",
			address:  0,
			err:      ErrSymbolAlreadyExists,
		},
	}

	for _, d := range data {
		d := d
		t.Run(d.testCase, func(t *testing.T) {
			t.Parallel()

			err := table.AddEntry(d.symbol, d.address)

			if d.err != nil {
				if diff := cmp.Diff(err, d.err, cmpopts.EquateErrors()); diff != "" {
					t.Fatal(diff)
				}
				return
			} else if err != nil {
				t.Fatal(err)
			}

			if !table.Contains(d.symbol) {
				t.Fatal("symbol not found")
			}

			address, err := table.GetAddress(d.symbol)
			if err != nil {
				t.Fatal(err)
			}
			if address != d.address {
				t.Fatalf("expected address %d, got %d", d.address, address)
			}
		})
	}
}
