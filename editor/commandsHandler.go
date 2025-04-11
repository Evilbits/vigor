package editor

import (
	"fmt"
	"os"

	"github.com/evilbits/vigor/ui"
	"github.com/gdamore/tcell/v2"
)

func (ed *Editor) HandleKey(event *tcell.EventKey) {
	fa, err := ed.screen.GetFocusedEditableArea()
	if err != nil {
		panic(err)
	}
	switch typedFa := fa.(type) {
	case *ui.TextArea:
		switch typedFa.Mode {
		case ui.VisualMode:
			if ed.cmd.CommandMode {
				ed.handleCmdCommandEvent(event)
			} else {
				ed.handleVisualModeKey(typedFa, event.Rune())
			}
		case ui.InsertMode:
			ed.handleInsertModeEvent(typedFa, event)
		}
	case *ui.FileBrowser:
		if ed.cmd.CommandMode {
			ed.handleCmdCommandEvent(event)
		} else {
			ed.handleFileBrowserKey(typedFa, event)
		}
	}
}

func (ed *Editor) handleFileBrowserKey(fb *ui.FileBrowser, event *tcell.EventKey) {
	switch event.Key() {
	case tcell.KeyEnter:
		browserFile := fb.GetCurrentFile()
		file := NewFile(browserFile.Name())
		ed.LoadFile(file)
		ed.screen.Grid.ReplaceCurrentFocusedEditableArea(ed.textArea)
	default:
		switch event.Rune() {
		case 'j':
			fb.MoveCursor(1)
		case 'k':
			fb.MoveCursor(-1)
		case ':':
			ed.cmd.StartCommandMode()
		}
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

func (ed *Editor) executeCmdCommand(command string) {
	ed.cmd.ExitCommandMode()
	switch command {
	case "q", "quit":
		ed.handleExit()
	case "w", "write":
		fa, err := ed.screen.GetFocusedEditableArea()
		if err != nil {
			panic(err)
		}

		if ta, ok := fa.(*ui.TextArea); ok {
			err = ed.activeFile.WriteFile(ta.TextContent)
			if err != nil {
				panic(err)
			}
			ed.cmd.AddText(fmt.Sprintf("Successfully wrote to: %v", ed.activeFile.absPath))
		} else {
			ed.cmd.SetError("Tried writing to a focused area that doesn't support writing")
		}
	case "e", "edit":
		// Replace currently active TextArea with a FileBrowser
		ed.screen.Grid.ReplaceCurrentFocusedEditableArea(ed.fileBrowser)
	default:
		ed.cmd.SetError(fmt.Sprintf("Invalid command: %s", command))
	}
}

func (ed *Editor) handleExit() {
	ed.screen.Fini()
	os.Exit(0)
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
