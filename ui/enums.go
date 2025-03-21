package ui

// LineFeed is a type that represents a line feed character
type LineFeed rune

// Define the constant
const (
	LF LineFeed = '\n'
)

// String returns the string representation of LineFeed
func (lf LineFeed) String() string {
	return string(lf)
}

// Add returns a string with the LineFeed appended
func (lf LineFeed) Add(s string) string {
	return s + string(lf)
}
