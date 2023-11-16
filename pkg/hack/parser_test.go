package hack

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	// Assemble code for testing.
	testAsm = `
  // comment

  // comment

  @i
  M=D+A;JMP
  (LOOP)
  `
)

func TestNewParser(t *testing.T) {
	t.Parallel()

	p := NewParser(strings.NewReader(testAsm))
	if p.r == nil {
		t.Fatal()
	}

	if p.HasMoreCommands() != true {
		t.Error("HasMoreCommands() should be true")
	}
}

func TestParser_Advance(t *testing.T) {
	t.Parallel()

	p := NewParser(strings.NewReader(testAsm))
	if p.Advance() != true {
		t.Error("Advance() should be true")
	}
	if p.Advance() != true {
		t.Error("Advance() should be true")
	}
	if p.Advance() != true {
		t.Error("Advance() should be true")
	}
	if p.Advance() != false {
		t.Error("Advance() should be false")
	}
}

func TestParser_Symbol_ACommand(t *testing.T) {
	t.Parallel()

	p := NewParser(strings.NewReader(testAsm))
	p.Advance()

	sym, err := p.Symbol()
	if err != nil {
		t.Fatal(err)
	}

	if sym != "i" {
		t.Error("Symbol() should return i")
	}
}

func TestParser_Symobol_LCommand(t *testing.T) {
	t.Parallel()

	p := NewParser(strings.NewReader(testAsm))
	p.Advance()
	p.Advance()
	p.Advance()

	sym, err := p.Symbol()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(sym, "LOOP"); diff != "" {
		t.Error(diff)
	}
}

func TestParser_CommandType(t *testing.T) {
	t.Parallel()

	p := NewParser(strings.NewReader(testAsm))
	p.Advance()

	if diff := cmp.Diff(p.CommandType(), ACommand); diff != "" {
		t.Error(diff)
	}

	p.Advance()
	if diff := cmp.Diff(p.CommandType(), CCommand); diff != "" {
		t.Error(diff)
	}

	p.Advance()
	if diff := cmp.Diff(p.CommandType(), LCommand); diff != "" {
		t.Error(diff)
	}
}

func TestParser_Dest(t *testing.T) {
	t.Parallel()

	p := NewParser(strings.NewReader(testAsm))
	p.Advance()
	p.Advance()

	dest, err := p.Dest()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(dest, "M"); diff != "" {
		t.Error(diff)
	}
}

func TestParser_Comp(t *testing.T) {
	t.Parallel()

	p := NewParser(strings.NewReader(testAsm))
	p.Advance()
	p.Advance()

	comp, err := p.Comp()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(comp, "D+A"); diff != "" {
		t.Error(diff)
	}
}

func TestParser_Jump(t *testing.T) {
	t.Parallel()

	p := NewParser(strings.NewReader(testAsm))
	p.Advance()
	p.Advance()

	jump, err := p.Jump()
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(jump, "JMP"); diff != "" {
		t.Error(diff)
	}
}
