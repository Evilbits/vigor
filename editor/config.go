package editor

import "github.com/evilbits/vigor/ui"

// TODO: Implement actual config reading
func ReadConf() *ui.Grid {
	rootGrid := ui.NewGrid()
	rootBox := ui.NewTextArea()

	rootBoxTwo := ui.NewBox()
	rootBoxTwo.SetBackgroundColor("grey")
	rootBoxTwo.AddText("Position")

	rootGrid.
		SetRows(0, 1).
		AddItem(rootBox, true).
		AddItem(rootBoxTwo, false)
	return rootGrid
}
