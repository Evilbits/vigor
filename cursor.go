package main

import (
	"fmt"
	"github.com/rivo/tview"
)

func getCursorPosStr(area *tview.TextArea) string {
	fromRow, fromColumn, toRow, toColumn := area.GetCursor()

	if fromRow == toRow && fromColumn == toColumn {
		return fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn)
	} else {
		return fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn)
	}
}
