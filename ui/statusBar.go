package ui

import "fmt"

type StatusBar struct {
	*Box
	*TextArea

	ActiveFileName string
}

func NewStatusBar(ta *TextArea) *StatusBar {
	sb := &StatusBar{}
	sb.Box = NewBox()
	sb.TextArea = ta
	return sb
}

func statusBarSeparator() string {
	return " | "
}

func (sb *StatusBar) Draw(screen *Screen) {
	statusBarStr := ""
	sb.Box.Draw(screen)
	statusBarStr += fmt.Sprintf("(%v, %v)", sb.TextArea.cursorX, sb.TextArea.cursorY)

	if sb.ActiveFileName != "" {
		statusBarStr += statusBarSeparator()
		statusBarStr += sb.ActiveFileName
	}

	sb.AddText(statusBarStr)
}

func (sb *StatusBar) AddText(text string) {
	sb.Text = text
}

func (sb *StatusBar) HandleKey(char rune, screen *Screen) {
	// Not implemented
}
