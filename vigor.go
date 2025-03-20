package main

import (
	"github.com/rivo/tview"
)

func buildTextArea() (*tview.TextArea, *tview.Box) {
	textArea := tview.NewTextArea().SetPlaceholder("Enter text here")
	textAreaBox := textArea.SetTitle("Vigor").
		SetBorder(true)

	return textArea, textAreaBox
}

func buildPosition(textArea *tview.TextArea) *tview.TextView {
	position := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignRight)

	updatePosFunc := func() {
		position.SetText(getCursorPosStr(textArea))
	}
	textArea.SetMovedFunc(updatePosFunc)
	updatePosFunc()

	return position
}

func buildStateArea() *tview.TextView {
	state := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("Visual")

	return state
}

func buildMainView() *tview.Grid {
	textArea, textAreaBox := buildTextArea()
	stateArea := buildStateArea()
	positionArea := buildPosition(textArea)

	mainLayout := tview.NewGrid().
		SetRows(0, 1).
		AddItem(textAreaBox, 0, 0, 1, 2, 0, 0, true).
		AddItem(stateArea, 1, 0, 1, 1, 0, 0, false).
		AddItem(positionArea, 1, 1, 1, 1, 0, 0, false)
	return mainLayout
}

func main() {
	app := tview.NewApplication()
	mainView := buildMainView()

	if err := app.SetRoot(mainView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
