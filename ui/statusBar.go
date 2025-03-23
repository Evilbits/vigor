package ui

import (
	"fmt"
)

type StatusBar struct {
	*Box

	monitoredTextArea *TextArea
	ActiveFileName    string
}

func NewStatusBar(ta *TextArea) *StatusBar {
	sb := &StatusBar{}
	sb.Box = NewBox()
	sb.monitoredTextArea = ta
	return sb
}

func statusBarSeparator() string {
	return " | "
}

func (sb *StatusBar) Draw(screen *Screen) {
	ta := sb.monitoredTextArea
	statusBarStr := ""
	statusBarStr += fmt.Sprintf("(%v [%v], %v, {%v})", ta.cursorX, ta.lastUserXPos, ta.cursorY, ta.getCursorLocInText())

	statusBarStr += statusBarSeparator()
	statusBarStr += fmt.Sprintf("[%v]", ta.Mode)

	if sb.ActiveFileName != "" {
		statusBarStr += statusBarSeparator()
		statusBarStr += sb.ActiveFileName
	}
	statusBarStr += statusBarSeparator()
	statusBarStr += fmt.Sprintf("%v", ta.GetTextContentOffset())

	sb.AddText(statusBarStr)
	sb.Box.Draw(screen)
}
