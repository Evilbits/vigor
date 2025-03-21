package editor

import (
	"os"

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

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data[:])
}

func (editor *Editor) Start(filepath string) {
	text := readFile(filepath)
	rootGrid := ui.NewGrid()
	rootBox := ui.NewBox(editor.screen)
	rootBox.AddText(text)

	rootBoxTwo := ui.NewBox(editor.screen)
	rootBoxTwo.SetBackgroundColor("red")
	rootBoxTwo.AddText("Position")

	rootGrid.
		SetRows(0, 1).
		AddItem(rootBox).
		AddItem(rootBoxTwo).
		Draw(editor.screen)

	editor.screen.StartEventLoop(rootGrid)
}
