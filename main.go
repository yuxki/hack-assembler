package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yuxki/hack-assembler/pkg/hack"
)

func printUsage() {
	fmt.Println("Usage: hack-assembler <asm file>")
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	asmFile := os.Args[1]

	if !strings.HasSuffix(asmFile, ".asm") {
		fmt.Printf("Error: file must have .asm extension: %s\n", asmFile)
		return
	}

	if _, err := os.Stat(asmFile); os.IsNotExist(err) {
		fmt.Printf("Error: asm file does not exist: %s\n", asmFile)
		return
	}

	reader, err := os.Open(asmFile)
	if err != nil {
		fmt.Printf("Error: could not open asm file: %s\n", err.Error())
	}
	defer reader.Close()

	outFile := strings.Replace(asmFile, ".asm", ".hack", 1)
	writer, err := os.Create(outFile)
	if err != nil {
		fmt.Printf("Error: could not create hack file: %s\n", err.Error())
	}
	defer writer.Close()

	assmbler, err := hack.NewAssembler(reader, writer)
	if err != nil {
		panic("Error: could not create assembler: " + err.Error())
	}

	err = assmbler.Assemble()
	if err != nil {
		fmt.Printf("Error: could not assemble file: %s\n", err.Error())
	}
}
