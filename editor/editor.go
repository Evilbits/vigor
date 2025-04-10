package editor

import (
	"os"

	"github.com/evilbits/vigor/ui"
)

type Editor struct {
	screen *ui.Screen

	activeFile  *ViFile
	textArea    *ui.TextArea
	fileBrowser *ui.FileBrowser
	statusBar   *ui.StatusBar

	cmd *ui.Cmd
}

func NewEditor() *Editor {
	editor := &Editor{}
	editor.screen = ui.NewScreen()
	editor.screen.OnKeyPress = editor.HandleKey
	return editor
}

func currentDir() (string, error) {
	return os.Getwd()
}

func (editor *Editor) LoadFile(entry os.DirEntry) {
	file := NewFile(entry.Name())
	text := file.ReadFileContents()

	editor.cmd.AddText(file.absPath)

	editor.textArea.TextContent = text
	editor.statusBar.ActiveFileName = file.GetFileName()
}

func (editor *Editor) Start(filepath string, debug bool) {
	editor.activeFile = NewFile(filepath)
	text := editor.activeFile.ReadFileContents()
	textArea, cmd, statusBar, grid := ReadConf(debug)
	editor.cmd = cmd
	editor.screen.Grid = grid
	editor.statusBar = statusBar

	textArea.TextContent = text
	statusBar.ActiveFileName = editor.activeFile.GetFileName()

	editor.textArea = textArea
	currDir, err := currentDir()
	if err != nil {
		panic("Could not resolve current directory")
	}
	editor.fileBrowser = ui.NewFileBrowser(currDir)

	editor.screen.StartEventLoop(textArea)
}
