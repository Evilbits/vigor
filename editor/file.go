package editor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	// "github.com/evilbits/vigor/ui"
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
	scanner := bufio.NewScanner(f)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
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
