package editor

import (
	"fmt"
	"os"
	"strings"

	"github.com/evilbits/vigor/ui"
)

type Editor struct {
	screen *ui.Screen

	activeFile  *File
	textArea    *ui.TextArea
	fileBrowser *ui.FileBrowser

	cmd *ui.Cmd
}

func NewEditor() *Editor {
	editor := &Editor{}
	editor.screen = ui.NewScreen()
	editor.screen.OnKeyPress = editor.HandleKey
	return editor
}

func filePathToFileName(filepath string) string {
	if strings.Contains(filepath, "/") {
		splitStr := strings.Split(filepath, "/")
		return fmt.Sprint(splitStr[len(splitStr)-1])
	}
	return filepath
}

func currentDir() (string, error) {
	return os.Getwd()
}

func (editor *Editor) Start(filepath string, debug bool) {
	editor.activeFile = NewFile(filepath)
	text := editor.activeFile.ReadFile()
	textArea, cmd, statusBar, grid := ReadConf(debug)
	editor.cmd = cmd
	editor.screen.Grid = grid

	textArea.TextContent = text
	statusBar.ActiveFileName = filePathToFileName(filepath)

	editor.textArea = textArea
	currDir, err := currentDir()
	if err != nil {
		panic("Could not resolve current directory")
	}
	editor.fileBrowser = ui.NewFileBrowser(currDir)

	editor.screen.StartEventLoop(textArea)
}
