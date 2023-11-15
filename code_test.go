package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCode_Dest(t *testing.T) {
	t.Parallel()

	c := NewCode()
	bits, err := c.Dest("M")
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(bits, "001"); diff != "" {
		t.Error(diff)
	}
}

func TestCode_Comp(t *testing.T) {
	t.Parallel()

	c := NewCode()
	bits, err := c.Comp("D+1")
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(bits, "0011111"); diff != "" {
		t.Error(diff)
	}
}

func TestCode_Jump(t *testing.T) {
	t.Parallel()

	c := NewCode()
	bits, err := c.Jump("JGT")
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(bits, "001"); diff != "" {
		t.Error(diff)
	}
}
