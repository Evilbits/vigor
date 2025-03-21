package editor

import (
	"github.com/evilbits/vigor/ui"
)

// TODO: Implement actual config reading
func ReadConf() *ui.Grid {
	rootGrid := ui.NewGrid()
	textArea := ui.NewTextArea()

	statusBar := ui.NewStatusBar(textArea)
	statusBar.SetBackgroundColor("grey")

	rootGrid.
		SetRows(0, 1).
		AddItem(textArea, true).
		AddItem(statusBar, false)
	return rootGrid
}
