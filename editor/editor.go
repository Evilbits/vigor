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

	activeFile *File

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

func (editor *Editor) Start(filepath string, debug bool) {
	editor.activeFile = NewFile(filepath)
	text := editor.activeFile.ReadFile()
	textArea, cmd, statusBar, grid := ReadConf(debug)
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
			ed.handleCmdCommandEvent(event)
		} else {
			ed.handleVisualModeKey(ta, char)
		}
	case ui.InsertMode:
		ed.handleInsertModeEvent(ta, event)
	}
}

func (ed *Editor) handleInsertModeEvent(ta *ui.TextArea, event *tcell.EventKey) {
	char := event.Rune()

	switch event.Key() {
	case tcell.KeyEsc:
		ta.Mode = ui.VisualMode
	case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyDelete:
		ta.RemoveChar()
	case tcell.KeyEnter:
		ta.InsertNewline()
	default:
		ta.InsertChar(char)
	}
}

func (ed *Editor) handleCmdCommandEvent(event *tcell.EventKey) {
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
	case "w", "write":
		ta, err := ed.screen.GetFocusedEditableArea()
		if err != nil {
			panic(err)
		}
		err = ed.activeFile.WriteFile(ta.TextContent)
		if err != nil {
			panic(err)
		}
		ed.cmd.AddText(fmt.Sprintf("Successfully wrote to: %v", ed.activeFile.path))
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
