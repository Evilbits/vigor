package editor

import (
	"github.com/evilbits/vigor/ui"
)

// TODO: Implement actual config reading
func ReadConf() (*ui.TextArea, *ui.Cmd, *ui.StatusBar, *ui.Grid) {
	rootGrid := ui.NewGrid()
	textArea := ui.NewTextArea()

	cmd := ui.NewCmd()

	statusBar := ui.NewStatusBar(textArea)
	statusBar.SetBackgroundColor("grey")

	rootGrid.
		SetRows(0, 1, 1).
		AddItem(textArea, true).
		AddItem(cmd, false).
		AddItem(statusBar, false)
	return textArea, cmd, statusBar, rootGrid
}
