package editor

import (
	"fmt"
	"os"
	"strings"

	"github.com/evilbits/vigor/ui"
)

type File struct {
	path string
}

func NewFile(path string) *File {
	file := &File{
		path: path,
	}
	return file
}

func (f *File) ReadFile() []string {
	data, err := os.ReadFile(f.path)
	if err != nil {
		panic(err)
	}

	return strings.Split(string(data[:]), fmt.Sprint(ui.LF))
}

func (f *File) WriteFile(text []string) error {
	err := assertFileExists(f.path)
	if err != nil {
		panic(err)
	}

	data := strings.Join(text, "\n")

	return os.WriteFile(f.path, []byte(data), 0644)
}

func assertFileExists(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	return nil
}
