package editor

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/evilbits/vigor/ui"
	"github.com/gdamore/tcell/v2"
)

type Editor struct {
	screen *ui.Screen

	cmd *ui.Cmd
}

func NewEditor() *Editor {
	editor := &Editor{}
	editor.screen = ui.NewScreen()
	editor.screen.OnKeyPress = editor.HandleKey
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
	textArea, cmd, statusBar, grid := ReadConf()
	editor.cmd = cmd
	editor.screen.Grid = grid

	textArea.TextContent = text
	statusBar.ActiveFileName = filePathToFileName(filepath)

	focusedArea, err := editor.screen.Grid.GetFocusedEditableArea()
	if err != nil {
		log.Fatal(err)
	}

	editor.screen.StartEventLoop(focusedArea)
}

func (ed *Editor) HandleKey(event *tcell.EventKey) {
	char := event.Rune()

	ta, err := ed.screen.GetFocusedEditableArea()
	if err != nil {
		panic(err)
	}
	ta.LastKeySeen = event.Key()
	switch ta.Mode {
	case ui.VisualMode:
		if ed.cmd.CommandMode {
			ed.handleCmdCommandKey(event)
		} else {
			ed.handleVisualModeKey(ta, char)
		}
	case ui.InsertMode:
		if event.Key() == tcell.KeyEsc {
			ta.Mode = ui.VisualMode
			return
		}
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 || event.Key() == tcell.KeyDelete {
			ta.RemoveChar()
			return
		}
		ta.InsertChar(char)
		ta.MoveCursor(1, 0)
	}
}

func (ed *Editor) handleCmdCommandKey(event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyEsc:
		ed.cmd.ExitCommandMode()
	case tcell.KeyEnter:
		ed.executeCmdCommand(ed.cmd.CurrentCommand)
	case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyDelete:
		ed.cmd.DeleteLastCharFromCommand()
	default:
		char := event.Rune()
		ed.cmd.AppendRuneToCurrentCommand(char)
	}
}

func (ed *Editor) executeCmdCommand(command string) {
	ed.cmd.ExitCommandMode()
	switch command {
	case "q", "quit":
		ed.handleExit()
	default:
		ed.cmd.SetError(fmt.Sprintf("Invalid command: %s", command))
	}
}

func (ed *Editor) handleVisualModeKey(ta *ui.TextArea, char rune) {
	switch char {
	case 'h':
		ta.MoveCursor(-1, 0)
	case 'j':
		ta.MoveCursor(0, 1)
	case 'k':
		ta.MoveCursor(0, -1)
	case 'l':
		ta.MoveCursor(1, 0)
	case 'i':
		ta.Mode = ui.InsertMode
	case 'a':
		ta.Mode = ui.InsertMode
		ta.MoveCursor(1, 0)
	case '$':
		ta.MoveCursorEndOfCurrLine()
	case '0':
		ta.MoveCursorBeginningOfCurrLine()
	case 'g':
		ta.MoveCursorBeginningOfFile()
	case 'G':
		ta.MoveCursorEndOfFile()
	case ':':
		ed.cmd.StartCommandMode()
	}
}

func (ed *Editor) handleExit() {
	ed.screen.Fini()
	os.Exit(0)
}
