package hack

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func removeComment(line string) string {
	return strings.Split(line, "//")[0]
}

// This is a parser for the Hack assembly language.
// - It parses the assembly language into its individual components.
// - It removes comments and whitespace.
// - It tracks the current line number.
// - It also keeps track of the current command type.
type Parser struct {
	r               io.Reader
	scanner         *bufio.Scanner
	hasMoreCommands bool
	lineNumber      uint

	regComment   *regexp.Regexp
	regSpaceLine *regexp.Regexp

	regACommandSymbol *regexp.Regexp

	regLCommandSymbol *regexp.Regexp

	regDest *regexp.Regexp
	regComp *regexp.Regexp
	regJump *regexp.Regexp
}

type CommandType int

const (
	ACommand CommandType = iota
	CCommand
	LCommand
)

// NewParser creates a new parser.
func NewParser(r io.Reader) *Parser {
	var p Parser

	p.Reset(r)

	p.regComment = regexp.MustCompile(`//.*$`)
	p.regSpaceLine = regexp.MustCompile(`^\s*$`)

	p.regACommandSymbol = regexp.MustCompile(`^\s*@([0-9a-z-A-Z_.$:][0-9a-z-A-Z_.$:]*)`)
	p.regLCommandSymbol = regexp.MustCompile(`^\s*\(([a-z-A-Z_.$:][0-9a-z-A-Z_.$:]*)\)`)

	p.regDest = regexp.MustCompile(`^\s*(M|D|MD|A|AM|AD|AMD)=`)
	p.regComp = regexp.MustCompile(
		`^(D\+1|A\+1|D-1|A-1|D\+A|D-A|A-D|D&A|D\|A|M\+1|M-1|D\+M|D-M|M-D|D\|M|D&M|0|1|-1|D|A|!D|!A|-D|-A|M|M|!M|-M)$`,
	)
	p.regJump = regexp.MustCompile(`;(JGT|JEQ|JGE|JLT|JNE|JLE|JMP)`)

	return &p
}

// Advance advances the parser to the next command
// Returns true if there are more commands to parse.
func (p *Parser) Advance() bool {
	for p.scanner.Scan() {
		line := p.scanner.Text()
		line = p.regComment.ReplaceAllString(line, "")
		if p.regSpaceLine.MatchString(line) {
			continue
		}

		if p.CommandType() != LCommand {
			p.lineNumber++
		}

		return true
	}
	p.hasMoreCommands = false
	return false
}

// HasMoreCommands returns true if there are more commands to parse.
func (p *Parser) HasMoreCommands() bool {
	return p.hasMoreCommands
}

// Command returns the current command.
func (p *Parser) Command() string {
	return p.scanner.Text()
}

// CommandType returns the type of the current command.
func (p *Parser) CommandType() CommandType {
	if regexp.MustCompile(`^\s*@`).MatchString(p.scanner.Text()) {
		return ACommand
	}
	if regexp.MustCompile(`^\s*\(`).MatchString(p.scanner.Text()) {
		return LCommand
	}
	return CCommand
}

// ErrNonAorLCommand is returned when the symbol command is called on a non-A or non-L command.
var ErrNonAorLCommand = errors.New("Dest called on non L or A command")

// Symbol returns the symbol of the current A or L command.
func (p *Parser) Symbol() (string, error) {
	if p.CommandType() == ACommand {
		if p.regACommandSymbol.MatchString(p.scanner.Text()) {
			return p.regACommandSymbol.FindStringSubmatch(p.scanner.Text())[1], nil
		}
	}

	if p.CommandType() == LCommand {
		if p.regLCommandSymbol.MatchString(p.scanner.Text()) {
			return p.regLCommandSymbol.FindStringSubmatch(p.scanner.Text())[1], nil
		}
	}

	return "", ErrNonAorLCommand
}

var ErrNonCCommand = errors.New("Dest called on non-C command")

// Dest returns the dest command of the current C command.
func (p *Parser) Dest() (string, error) {
	if p.CommandType() != CCommand {
		return "", ErrNonCCommand
	}

	if p.regDest.MatchString(p.scanner.Text()) {
		return p.regDest.FindStringSubmatch(p.scanner.Text())[1], nil
	}

	return "", nil
}

var ErrInvalidCompCommand = errors.New("invalid comp")

// Comp returns the comp command of the current C command.
func (p *Parser) Comp() (string, error) {
	if p.CommandType() != CCommand {
		return "", ErrNonCCommand
	}

	command := removeComment(p.scanner.Text())
	command = p.regDest.ReplaceAllString(command, "")
	command = p.regJump.ReplaceAllString(command, "")
	command = strings.TrimSpace(command)

	if p.regComp.MatchString(command) {
		return p.regComp.FindStringSubmatch(command)[1], nil
	}

	return "", fmt.Errorf("%s: %w", command, ErrInvalidCompCommand)
}

// Jump returns the jump command of the current C command.
func (p *Parser) Jump() (string, error) {
	if p.CommandType() != CCommand {
		return "", ErrNonCCommand
	}

	if p.regJump.MatchString(p.scanner.Text()) {
		return p.regJump.FindStringSubmatch(removeComment(p.scanner.Text()))[1], nil
	}

	return "", nil
}

// LineNumber returns the current line number of the parser.
func (p *Parser) LineNumber() uint {
	return p.lineNumber
}

// Reset resets the parser to read from the given reader.
func (p *Parser) Reset(r io.Reader) {
	p.r = r
	p.scanner = bufio.NewScanner(r)
	p.hasMoreCommands = true
	p.lineNumber = 0
}
