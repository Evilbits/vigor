package editor

import (
	"fmt"
	"github.com/evilbits/vigor/ui"
	"os"
	"path/filepath"
	"strings"
)

type ViFile struct {
	*os.File
	absPath string
}

func NewFile(path string) *ViFile {
	file, err := readFile(path)
	if err != nil {
		panic("Tried opening non existant file")
	}
	absPath, err := filepath.Abs(file.Name())
	if err != nil {
		panic("Could not resolve absolute path for file")
	}

	viFile := &ViFile{
		File:    file,
		absPath: absPath,
	}
	return viFile
}

func readFile(path string) (*os.File, error) {
	err := assertFileExists(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *ViFile) ReadFileContents() []string {
	content, err := os.ReadFile(f.absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return nil
	}

	return strings.Split(string(content), fmt.Sprint(ui.LF))
}

func (f *ViFile) WriteFile(text []string) error {
	err := assertFileExists(f.absPath)
	if err != nil {
		panic(err)
	}

	data := strings.Join(text, "\n")

	return os.WriteFile(f.absPath, []byte(data), 0644)
}

func (f *ViFile) GetFileName() string {
	if strings.Contains(f.absPath, "/") {
		splitStr := strings.Split(f.absPath, "/")
		return fmt.Sprint(splitStr[len(splitStr)-1])
	}
	return f.absPath
}

func assertFileExists(filepath string) error {
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	return nil
}
