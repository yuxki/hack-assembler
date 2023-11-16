package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuxki/hack-assembler/pkg/hack"
)

func printUsage() {
	println("Usage: hack-assembler <asm file>")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	if !strings.HasSuffix(os.Args[1], ".asm") {
		println("Error: file must have .asm extension")
		return
	}

	asmFile := os.Args[1]
	if _, err := os.Stat(asmFile); os.IsNotExist(err) {
		fmt.Printf("Error: asm file does not exist\n", err.Error())
		return
	}

	reader, err := os.Open(asmFile)
	if err != nil {
		fmt.Printf("Error: could not open asm file: %s\n", err.Error())
	}
	defer reader.Close()

	parser := hack.NewParser(reader)
	code := hack.NewCode()
	symbolTable, err := hack.NewSymbolTable()
	if err != nil {
		panic("Error: could not create symbol table.")
	}

	outFile := strings.Replace(asmFile, ".asm", ".hack", 1)
	writer, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("Error: could not create hack file: %s\n", err.Error())
	}
	defer writer.Close()

	assmbler := hack.NewAssembler(writer, parser, code, symbolTable)
	err = assmbler.Assemble()
	if err != nil {
		fmt.Printf("Error: could not assemble file: %s\n", err.Error())
	}
}
