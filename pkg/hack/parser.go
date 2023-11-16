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

func NewParser(r io.Reader) *Parser {
	var p Parser

	p.Reset(r)

	p.regComment = regexp.MustCompile(`//.*$`)
	p.regSpaceLine = regexp.MustCompile(`^\s*$`)

	p.regACommandSymbol = regexp.MustCompile(`^\s*@([0-9a-z-A-Z_.$:][0-9a-z-A-Z_.$:]*)`)
	p.regLCommandSymbol = regexp.MustCompile(`^\s*\(([a-z-A-Z_.$:][0-9a-z-A-Z_.$:]*)\)`)

	p.regDest = regexp.MustCompile(`^\s*(M|D|MD|A|AM|AD|AMD)=`)
	p.regComp = regexp.MustCompile(`^(D\+1|A\+1|D-1|A-1|D\+A|D-A|A-D|D&A|D\|A|M\+1|M-1|D\+M|D-M|M-D|D\|M|D&M|0|1|-1|D|A|!D|!A|-D|-A|M|M|!M|-M)$`)
	p.regJump = regexp.MustCompile(`;(JGT|JEQ|JGE|JLT|JNE|JLE|JMP)`)

	return &p
}

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

func (p *Parser) HasMoreCommands() bool {
	return p.hasMoreCommands
}

func (p *Parser) Command() string {
	return p.scanner.Text()
}

func (p *Parser) CommandType() CommandType {
	if regexp.MustCompile(`^\s*@`).MatchString(p.scanner.Text()) {
		return ACommand
	}
	if regexp.MustCompile(`^\s*\(`).MatchString(p.scanner.Text()) {
		return LCommand
	}
	return CCommand
}

var ErrNonAorLCommand = errors.New("Dest called on non L or A command")

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

func (p *Parser) Dest() (string, error) {
	if p.CommandType() != CCommand {
		return "", ErrNonCCommand
	}

	if p.regDest.MatchString(p.scanner.Text()) {
		return p.regDest.FindStringSubmatch(p.scanner.Text())[1], nil
	}

	return "", nil
}

var ErrInvalidCompCommand = errors.New("Invalid comp")

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

func (p *Parser) Jump() (string, error) {
	if p.CommandType() != CCommand {
		return "", ErrNonCCommand
	}

	if p.regJump.MatchString(p.scanner.Text()) {
		return p.regJump.FindStringSubmatch(removeComment(p.scanner.Text()))[1], nil
	}

	return "", nil
}

func (p *Parser) LineNumber() uint {
	return p.lineNumber
}

func (p *Parser) Reset(r io.Reader) {
	p.r = r
	p.scanner = bufio.NewScanner(r)
	p.hasMoreCommands = true
	p.lineNumber = 0
}
