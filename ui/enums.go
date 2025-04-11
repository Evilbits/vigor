package ui

type LineFeed rune

const (
	LF LineFeed = '\n'
)

func (lf LineFeed) String() string {
	return string(lf)
}
