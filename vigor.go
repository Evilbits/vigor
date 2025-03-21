package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/evilbits/vigor/editor"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "f", "", "Path to the file to open")
	flag.Parse()

	if filePath == "" && len(flag.Args()) > 0 {
		filePath = flag.Args()[0]
	}

	filePath = fmt.Sprintf("./%s", filePath)
	log.Printf("Opening file: %s", filePath)

	editor.NewEditor().Start(filePath)
}
