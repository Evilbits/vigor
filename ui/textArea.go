package ui

import (
	"errors"
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
	// Stores last x pos that a user moved to. This allows better behaviour when going from a long line
	// to a short line and then to a long line again
	lastUserXPos int

	// Debug
	lastKeySeen tcell.Key
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

func (ta *TextArea) HandleKey(event *tcell.EventKey) {
	char := event.Rune()
	ta.lastKeySeen = event.Key()
	switch ta.mode {
	case VisualMode:
		ta.handleVisualModeKey(char)
	case InsertMode:
		if event.Key() == tcell.KeyEsc {
			ta.mode = VisualMode
			return
		}
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyBackspace2 || event.Key() == tcell.KeyDelete {
			ta.removeChar(ta.cursorX, ta.cursorY)
			return
		}
		ta.insertChar(ta.cursorX, ta.cursorY, char)
		ta.cursorX = ta.cursorX + 1
	}
}

func (ta *TextArea) handleVisualModeKey(char rune) {
	switch char {
	case 'h':
		ta.moveCursor(-1, 0)
	case 'j':
		ta.moveCursor(0, 1)
	case 'k':
		ta.moveCursor(0, -1)
	case 'l':
		ta.moveCursor(1, 0)
	case 'i':
		ta.mode = InsertMode
	case 'a':
		ta.mode = InsertMode
		ta.moveCursor(1, 0)
	}
}

func (ta *TextArea) SetMode(mode Mode) {
	ta.mode = mode
}

func (ta *TextArea) GetMode() Mode {
	return ta.mode
}

func (ta *TextArea) moveCursor(moveX int, moveY int) {
	x, y := ta.GetXY()
	if (ta.cursorX+moveX < x) || (ta.cursorY+moveY < y) {
		return
	}
	width, height := ta.Box.GetSize()
	if (ta.cursorX+moveX >= width) || (ta.cursorY+moveY >= height) {
		return
	}

	if moveX != 0 {
		ta.lastUserXPos = ta.cursorX + moveX
	}

	rowLen := len(ta.TextContent[ta.cursorY])
	// Don't allow moving further on the x axis than the content
	if ta.cursorX+moveX > rowLen {
		return
	}
	// moveY is required to be checked here as we can be at x+1 position after leaving insert mode
	if rowLen > 0 && ta.cursorX+moveX >= rowLen && ta.mode == VisualMode && moveY == 0 {
		return
	}

	// If we are y moving fit to the length of the new row on the x axis
	// This implements more IDE-like cursor movement between lines of different length
	if moveY != 0 {
		if ta.cursorY+moveY >= 0 && ta.cursorY+moveY < len(ta.TextContent) {
			nextRowLen := len(ta.TextContent[ta.cursorY+moveY])
			if ta.lastUserXPos > nextRowLen && ta.cursorX < nextRowLen {
				// Move to end of row if we have stored a prev x value that's higher
				moveX += nextRowLen - ta.cursorX - 1
			} else if ta.lastUserXPos <= nextRowLen {
				// Move to stored x val if next row is longer
				toMove := ta.lastUserXPos - ta.cursorX
				if ta.cursorX+toMove == nextRowLen && nextRowLen != 0 {
					// Handle a case where we move from rowLen+1 from another line
					toMove += -1
				}
				moveX = toMove
			} else if ta.cursorX >= nextRowLen {
				// If next row is shorter than our current row move to end of it
				if nextRowLen == 0 {
					moveX -= ta.cursorX - nextRowLen
				} else {
					moveX -= ta.cursorX - nextRowLen + 1
				}
			}
		}
	}

	// Don't allow going outside Y axis of text
	if ta.cursorY+moveY >= len(ta.TextContent) {
		return
	}

	ta.cursorX += moveX
	ta.cursorY += moveY
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
	ta.lastUserXPos += 1
	ta.TextContent[y] = currStr[:x] + string(char) + currStr[x:]
}

func (ta *TextArea) removeChar(x int, y int) error {
	currStr := ta.TextContent[y]
	if x == 0 || x > len(currStr) || y > len(ta.TextContent) {
		return errors.New("Cannot remove char out of bounds")
	}
	ta.TextContent[y] = currStr[:x-1] + currStr[x:]
	ta.moveCursor(-1, 0)
	return nil
}

func buildText(textArr []string) string {
	return strings.Join(textArr, fmt.Sprint(LF))
}
