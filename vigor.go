package main

import (
	"flag"

	"github.com/evilbits/vigor/editor"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "f", "", "Path to the file to open")
	flag.Parse()

	if filePath == "" && len(flag.Args()) > 0 {
		filePath = flag.Args()[0]
	}

	editor.NewEditor().Start(filePath)
}
