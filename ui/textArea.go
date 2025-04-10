package ui

import (
	"errors"
	"fmt"
	"slices"
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

	Mode Mode

	// Text content split by line delimiters
	TextContent []string
	// Allows y axis scroll of text ouside rendered content
	textContentOffset int

	cursorX, cursorY                     int
	cursorStyleVisual, cursorStyleInsert tcell.CursorStyle

	// Stores last x pos that a user moved to. This allows better behaviour when going from a long line
	// to a short line and then to a long line again
	lastUserXPos int
}

func NewTextArea() *TextArea {
	textArea := &TextArea{
		cursorX:           0,
		cursorY:           0,
		textContentOffset: 0,
		cursorStyleVisual: tcell.CursorStyleDefault,
		cursorStyleInsert: tcell.CursorStyleSteadyBar,
	}
	textArea.Box = NewBox()

	return textArea
}

func (ta *TextArea) MoveCursor(moveX int, moveY int) {
	x, y := ta.GetXY()
	// If going out of bounds
	if ta.cursorX+moveX < x {
		return
	}

	rowLen := len(ta.TextContent[ta.textContentOffset+ta.cursorY])
	if moveX != 0 {
		newUserXPos := min(ta.cursorX+moveX, rowLen-1)
		// If we don't check here newUserXPos will be -1 on an empty row
		if rowLen == 0 {
			newUserXPos += 1
		}
		ta.lastUserXPos = newUserXPos
	}

	// Don't allow moving further on the x axis than the content
	if ta.cursorX+moveX > rowLen {
		moveX = rowLen - 1 - ta.cursorX
	}
	// moveY is required to be checked here as we can be at x+1 position after leaving insert mode
	if rowLen > 0 && ta.cursorX+moveX >= rowLen && ta.Mode == VisualMode && moveY == 0 {
		return
	}

	textContentLen := len(ta.TextContent)
	cursorLoc := ta.getCursorLocInText()

	// If we are y moving fit to the length of the new row on the x axis
	// This implements more IDE-like cursor movement between lines of different length
	if moveY != 0 {
		nextRowIdx := cursorLoc + moveY

		if nextRowIdx >= 0 && nextRowIdx < textContentLen {
			nextRow := ta.TextContent[nextRowIdx]
			nextRowLen := len(nextRow)
			if ta.lastUserXPos > nextRowLen-1 {
				// Bound moveX by current x pos to prevent moving out of bounds
				moveX += max(nextRowLen-ta.cursorX, -ta.cursorX)
				if nextRowLen > 0 {
					moveX += -1
				}
			} else if ta.lastUserXPos < nextRowLen {
				toMove := ta.lastUserXPos - ta.cursorX
				if ta.cursorX+toMove == nextRowLen && nextRowLen != 0 {
					// Handle a case where we move from rowLen+1 from another line
					toMove += -1
				}
				moveX = toMove
			}
		}
	}

	_, renderedY := ta.Box.GetSize()
	// TODO: Clean up this and next piece of logic
	// Don't allow going outside Y axis of text. Simply bound the movement to the maximum rendered area that
	// we allow movement within
	if cursorLoc+moveY >= textContentLen {
		moveY = renderedY - ta.cursorY - 1
		ta.textContentOffset = textContentLen - renderedY
		ta.cursorY += moveY
		ta.cursorX += moveX
		return
	}

	// If we are scrolling outside rendered content in y axis but there is more text
	// scroll the screen and bound ypos to max allowed by bounding box
	if moveY != 0 {
		if moveY > 0 && ta.cursorY+moveY >= renderedY && textContentLen >= renderedY {
			// Scrolling down
			ta.textContentOffset += moveY
			if ta.textContentOffset+renderedY > textContentLen {
				ta.textContentOffset = textContentLen - renderedY
			}
			moveY = 0
		} else if moveY < 0 && ta.cursorY+moveY < y {
			// Scrolling up
			ta.textContentOffset = max(ta.textContentOffset+moveY, 0)
			moveY = 0
		}
	}

	ta.cursorX += moveX
	ta.cursorY += moveY

	if ta.cursorX < 0 {
		ta.cursorX = 0
	}
	if ta.cursorY < 0 {
		ta.cursorY = 0
	}
}

func (ta *TextArea) getCursorLocInText() int {
	return ta.cursorY + ta.textContentOffset
}

// Insert char at current cursor position
// Since TextContent is split by line we can use y as an index into our TextContent
func (ta *TextArea) InsertChar(char rune) {
	x, _ := ta.cursorX, ta.cursorY
	textY := ta.getCursorLocInText()
	currStr := ta.TextContent[textY]
	if x > len(currStr) {
		return
	}
	ta.TextContent[textY] = currStr[:x] + string(char) + currStr[x:]
	ta.MoveCursor(1, 0)
}

// Removes a char at current cursor position
func (ta *TextArea) RemoveChar() error {
	x := ta.cursorX
	textY := ta.getCursorLocInText()

	if x == 0 && textY != 0 {
		// If we are deleting from the beginning of our current line we need to merge with line above
		prevLineIdx := textY - 1
		prevLineLen := len(ta.TextContent[prevLineIdx])
		// Move first to make it easier to move to end of the line
		ta.MoveCursor(0, -1)
		ta.MoveCursorEndOfCurrLine()

		if prevLineLen > 0 {
			newLine := []string{ta.TextContent[prevLineIdx], ta.TextContent[textY]}
			ta.TextContent[prevLineIdx] = strings.Join(newLine, "")
			ta.DeleteLine(textY)
		} else {
			// Line is empty so we must delete it
			ta.DeleteLine(prevLineIdx)
		}
	} else {
		currStr := ta.TextContent[textY]
		if x == 0 || x > len(currStr) || textY > len(ta.TextContent) {
			return errors.New("Cannot remove char out of bounds")
		}

		ta.TextContent[textY] = currStr[:x-1] + currStr[x:]
		ta.MoveCursor(-1, 0)
	}
	return nil
}

func (ta *TextArea) DeleteLine(idx int) {
	ta.TextContent = slices.Delete(ta.TextContent, idx, idx+1)
}

func (ta *TextArea) InsertNewline() {
	currPos := ta.getCursorLocInText()
	currText := ta.TextContent[currPos]
	x := ta.cursorX
	strOne, strTwo := currText[:x], currText[x:]
	newTextContent := slices.Insert(ta.TextContent, currPos+1, strTwo)
	newTextContent[currPos] = strOne
	ta.TextContent = newTextContent
	ta.MoveCursor(-ta.cursorX, 1)
}

func (ta *TextArea) Draw(screen *Screen) {
	ta.Box.AddText(ta.buildTextContent())
	ta.Box.Draw(screen)

	cursorStyle := ta.cursorStyleVisual
	if ta.Mode == InsertMode {
		cursorStyle = ta.cursorStyleInsert
	}
	screen.RenderCursor(ta.cursorX, ta.cursorY, cursorStyle)
}

func (ta *TextArea) MoveCursorEndOfCurrLine() {
	var rowLen int
	if ta.Mode == VisualMode {
		rowLen = len(ta.TextContent[ta.textContentOffset+ta.cursorY]) - 1
	} else {
		rowLen = len(ta.TextContent[ta.textContentOffset+ta.cursorY])
	}

	moveX := rowLen - ta.cursorX
	ta.MoveCursor(moveX, 0)
}

func (ta *TextArea) MoveCursorBeginningOfCurrLine() {
	ta.MoveCursor(-ta.cursorX, 0)
}

func (ta *TextArea) MoveCursorBeginningOfFile() {
	ta.textContentOffset = 0
	ta.lastUserXPos = 0
	ta.cursorX = 0
	ta.cursorY = 0
}

func (ta *TextArea) MoveCursorEndOfFile() {
	ta.MoveCursor(-ta.cursorX, len(ta.TextContent)-ta.getCursorLocInText())
}

func (ta *TextArea) GetTextContentOffset() int {
	return ta.textContentOffset
}

func (ta *TextArea) buildTextContent() string {
	offset := ta.textContentOffset
	targetIdx := ta.getMaxAllowedYTargetOffset()
	textToRender := ta.TextContent[offset:targetIdx]

	return strings.Join(textToRender, fmt.Sprint(LF))
}

func (ta *TextArea) getMaxAllowedYTargetOffset() int {
	_, y := ta.Box.GetSize()
	offset := ta.textContentOffset
	targetIdx := min(y+offset, len(ta.TextContent))
	return targetIdx
}
