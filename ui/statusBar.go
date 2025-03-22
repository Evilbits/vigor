package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

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
	statusBarStr += fmt.Sprintf("(%v [%v], %v)", sb.TextArea.cursorX, sb.TextArea.lastUserXPos, sb.TextArea.cursorY)

	statusBarStr += statusBarSeparator()
	statusBarStr += fmt.Sprintf("[%v]", sb.GetMode())

	if sb.ActiveFileName != "" {
		statusBarStr += statusBarSeparator()
		statusBarStr += sb.ActiveFileName
	}

	sb.AddText(statusBarStr)
}

func (sb *StatusBar) AddText(text string) {
	sb.Text = text
}

// TODO: Find a way to remove
func (sb *StatusBar) HandleKey(event *tcell.EventKey, screen *Screen) {
	// Not implemented
}
