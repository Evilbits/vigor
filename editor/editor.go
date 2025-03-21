package editor

import (
	"os"

	"github.com/evilbits/vigor/ui"
)

type Editor struct {
	screen *ui.Screen
	grid   *ui.Grid
}

func NewEditor() *Editor {
	editor := &Editor{}
	editor.screen = ui.NewScreen()
	return editor
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data[:])
}

func (editor *Editor) Start(filepath string) {
	text := readFile(filepath)
	editor.grid = ReadConf()

	editor.grid.GetItem(0).Item.AddText(text)

	editor.grid.Draw(editor.screen)
	editor.screen.StartEventLoop(editor.grid)
}
