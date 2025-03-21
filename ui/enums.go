package ui

// LineFeed is a type that represents a line feed character
type LineFeed rune

const (
	LF LineFeed = '\n'
)

func (lf LineFeed) String() string {
	return string(lf)
}

func (lf LineFeed) Add(s string) string {
	return s + string(lf)
}
