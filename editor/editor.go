package editor

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/evilbits/vigor/ui"
)

type Editor struct {
	screen *ui.Screen
}

func NewEditor() *Editor {
	editor := &Editor{}
	editor.screen = ui.NewScreen()
	return editor
}

func readFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return strings.Split(string(data[:]), fmt.Sprint(ui.LF))
}

func filePathToFileName(filepath string) string {
	if strings.Contains(filepath, "/") {
		splitStr := strings.Split(filepath, "/")
		return fmt.Sprint(splitStr[len(splitStr)-1])
	}
	return filepath
}

func (editor *Editor) Start(filepath string) {
	text := readFile(filepath)
	grid := ReadConf()
	editor.screen.Grid = grid

	textArea, err := grid.GetFocusedEditableArea()
	if err != nil {
		log.Fatal(err)
	}
	textArea.TextContent = text

	statusBar, err := grid.GetStatusBar()
	if err != nil {
		log.Fatal(err)
	}
	statusBar.ActiveFileName = filePathToFileName(filepath)

	editor.screen.StartEventLoop()
}
