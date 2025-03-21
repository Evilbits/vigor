package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Mode int

func (m Mode) String() string {
	modes := [...]string{"Visual", "Insert"}
	if m < 0 || int(m) >= len(modes) {
		return "Unknown"
	}
	return modes[m]
}

const (
	VisualMode Mode = iota
	InsertMode
)

type TextArea struct {
	*Box

	mode Mode

	// Text content split by line delimiters
	TextContent      []string
	cursorX, cursorY int
}

func NewTextArea() *TextArea {
	textArea := &TextArea{
		cursorX: 0,
		cursorY: 0,
	}
	textArea.Box = NewBox()

	return textArea
}

func (ta *TextArea) Draw(screen *Screen) {
	ta.Box.AddText(buildText(ta.TextContent))
	ta.Box.Draw(screen)

	screen.RenderCursor(ta.cursorX, ta.cursorY)
}

func (ta *TextArea) HandleKey(event *tcell.EventKey, screen *Screen) {
	char := event.Rune()
	switch ta.mode {
	case VisualMode:
		ta.handleVisualModeKey(char, screen)
	case InsertMode:
		if event.Key() == tcell.KeyEsc && ta.mode == InsertMode {
			ta.mode = VisualMode
			return
		}
		ta.insertChar(ta.cursorX, ta.cursorY, char)
		ta.cursorX = ta.cursorX + 1
	}
}

func (ta *TextArea) handleVisualModeKey(char rune, screen *Screen) {
	switch char {
	case 'h':
		ta.moveCursor(-1, 0, screen)
	case 'j':
		ta.moveCursor(0, 1, screen)
	case 'k':
		ta.moveCursor(0, -1, screen)
	case 'l':
		ta.moveCursor(1, 0, screen)
	case 'i':
		ta.mode = InsertMode
	}
}

func (ta *TextArea) SetMode(mode Mode) {
	ta.mode = mode
}

func (ta *TextArea) GetMode() Mode {
	return ta.mode
}

func (ta *TextArea) moveCursor(moveX int, moveY int, screen *Screen) {
	x, y := ta.GetXY()
	if (ta.cursorX+moveX < x) || (ta.cursorY+moveY < y) {
		return
	}
	width, height := ta.Box.GetSize()
	if (ta.cursorX+moveX >= width) || (ta.cursorY+moveY >= height) {
		return
	}

	rowLen := len(ta.TextContent[ta.cursorY])
	// TODO: There is a lag of 1 row here. Fix
	// If we are y moving fit to the length of the new row on the x axis
	if moveY != 0 && ta.cursorX+moveX > rowLen {
		ta.cursorX = rowLen
	}

	// Don't allow moving further on the x axis than the content
	if ta.cursorX+moveX > rowLen {
		return
	}

	// Don't allow going outside Y axis of text
	if ta.cursorY+moveY >= len(ta.TextContent) {
		return
	}

	prevX := ta.cursorX
	prevY := ta.cursorY
	ta.cursorX += moveX
	ta.cursorY += moveY
	screen.RenderCursorMove(ta.cursorX, ta.cursorY, prevX, prevY)
}

// Since TextContent is split by line we can use y as an index into our TextContent
func (ta *TextArea) insertChar(x int, y int, char rune) {
	if y > len(ta.TextContent) {
		return
	}
	currStr := ta.TextContent[y]
	if x > len(currStr) {
		return
	}
	ta.TextContent[y] = currStr[:x] + string(char) + currStr[x:]
}

func buildText(textArr []string) string {
	return strings.Join(textArr, fmt.Sprint(LF))
}
