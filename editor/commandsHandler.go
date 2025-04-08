package editor

import (
	"github.com/evilbits/vigor/ui"
	"github.com/gdamore/tcell/v2"
)

func (ed *Editor) HandleKey(event *tcell.EventKey) {
	char := event.Rune()

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
				ed.handleVisualModeKey(typedFa, char)
			}
		case ui.InsertMode:
			ed.handleInsertModeEvent(typedFa, event)
		}
	case *ui.FileBrowser:
		if ed.cmd.CommandMode {
			ed.handleCmdCommandEvent(event)
		} else {
			ed.handleFileBrowserKey(typedFa, char)
		}
	}
}

func (ed *Editor) handleFileBrowserKey(fb *ui.FileBrowser, char rune) {
	switch char {
	case 'j':
		fb.MoveCursor(1)
	case 'k':
		fb.MoveCursor(-1)
	case ':':
		ed.cmd.StartCommandMode()
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
