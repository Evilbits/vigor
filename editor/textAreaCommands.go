package editor

import (
	"fmt"
	"github.com/evilbits/vigor/ui"
	"github.com/gdamore/tcell/v2"
	"os"
)

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
