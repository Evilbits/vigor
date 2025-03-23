package ui

import (
	"fmt"
)

type StatusBar struct {
	*Box

	monitoredTextArea *TextArea
	ActiveFileName    string
	DebugEnabled      bool
}

func NewStatusBar(ta *TextArea) *StatusBar {
	sb := &StatusBar{
		DebugEnabled: false,
	}
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
	if sb.DebugEnabled {
		statusBarStr += fmt.Sprintf("(%v [%v], %v, {%v})", ta.cursorX, ta.lastUserXPos, ta.cursorY, ta.getCursorLocInText())
	} else {
		statusBarStr += fmt.Sprintf("(%v , %v)", ta.cursorX, ta.cursorY)
	}

	statusBarStr += statusBarSeparator()
	statusBarStr += fmt.Sprintf("[%v]", ta.Mode)

	if sb.ActiveFileName != "" {
		statusBarStr += statusBarSeparator()
		statusBarStr += sb.ActiveFileName
	}
	if sb.DebugEnabled {
		statusBarStr += statusBarSeparator()
		statusBarStr += fmt.Sprintf("Offset: %v", ta.GetTextContentOffset())
	}

	sb.AddText(statusBarStr)
	sb.Box.Draw(screen)
}
