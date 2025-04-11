package ui

import (
	"errors"
	"fmt"
	"os"

	"slices"

	"github.com/gdamore/tcell/v2"
)

type FileNode struct {
	Path           string
	Name           string
	IsDir          bool
	IsOpen         bool
	subFiles       []*FileNode
	subDirectories []*FileNode
}

type FileBrowser struct {
	*Box

	rootNode *FileNode
	currDir  string

	cursorY int
}

func NewFileBrowser(currDir string) *FileBrowser {
	fileBrowser := &FileBrowser{
		currDir: currDir,
		cursorY: 0,
	}
	rootFile, err := fileBrowser.loadDir(currDir)
	if err != nil {
		panic(err)
	}
	fileBrowser.rootNode = rootFile
	fileBrowser.Box = NewBox()
	return fileBrowser
}

func (fb *FileBrowser) MoveCursor(moveY int) {
	fb.cursorY += moveY
	if fb.cursorY < 0 {
		fb.cursorY = 0
	}
	maxYPos := fb.countTreeNodes()
	if fb.cursorY > maxYPos {
		fb.cursorY = maxYPos
	}
}

func (fb *FileBrowser) GetCurrentNode() *FileNode {
	return fb.findNodeByCursorIdx()
}

func (fb *FileBrowser) loadDir(dir string) (*FileNode, error) {
	rootDir := &FileNode{
		Path:           dir,
		IsDir:          true,
		IsOpen:         false,
		subFiles:       []*FileNode{},
		subDirectories: []*FileNode{},
	}

	err := rootDir.addChildren()
	if err != nil {
		panic(err)
	}
	return rootDir, nil
}

func (fn *FileNode) addChildren() error {
	dirItems, err := os.ReadDir(fn.Path)
	if err != nil {
		return err
	}

	if !fn.IsOpen {
		for _, dirItem := range dirItems {
			node := &FileNode{
				Name:           dirItem.Name(),
				Path:           fmt.Sprintf("%v/%v", fn.Path, dirItem.Name()),
				IsDir:          dirItem.IsDir(),
				IsOpen:         false,
				subFiles:       []*FileNode{},
				subDirectories: []*FileNode{},
			}
			if node.IsDir {
				fn.subDirectories = append(fn.subDirectories, node)
			} else {
				fn.subFiles = append(fn.subFiles, node)
			}
		}
	} else {
		fn.subDirectories = []*FileNode{}
		fn.subFiles = []*FileNode{}
	}
	fn.IsOpen = !fn.IsOpen
	return nil
}

func (fb *FileBrowser) OpenDir(node *FileNode) error {
	if !node.IsDir {
		return errors.New("Tried to expand a file as a directory")
	}
	return node.addChildren()
}

func (fb *FileBrowser) Draw(screen *Screen) {
	fb.Box.AddText(fb.buildTextContent())
	fb.Box.Draw(screen)
	screen.RenderCursor(0, fb.cursorY, tcell.CursorStyleDefault)
}

func (fb *FileBrowser) buildTextContent() string {
	return fb.stringifyDirsAndFiles(fb.rootNode, "", 0)
}

func (fb *FileBrowser) findNodeByCursorIdx() *FileNode {
	targetIdx := fb.cursorY
	currentIdx := -1
	var result *FileNode

	var traverse func(*FileNode) bool
	traverse = func(node *FileNode) bool {
		if node != fb.rootNode {
			currentIdx++
		}
		if currentIdx == targetIdx {
			result = node
			return true
		}

		if slices.ContainsFunc(node.subDirectories, traverse) {
			return true
		}

		if slices.ContainsFunc(node.subFiles, traverse) {
			return true
		}
		return false
	}

	traverse(fb.rootNode)
	return result
}

func (fb *FileBrowser) countTreeNodes() int {
	return fb.walkAndCountTreeNodes(fb.rootNode, 0) - 1
}

func (fb *FileBrowser) walkAndCountTreeNodes(node *FileNode, count int) int {
	if node != fb.rootNode {
		count += 1
	}
	for _, dir := range node.subDirectories {
		count = fb.walkAndCountTreeNodes(dir, count)
	}
	for _, file := range node.subFiles {
		count = fb.walkAndCountTreeNodes(file, count)
	}

	return count
}

func (fb *FileBrowser) stringifyDirsAndFiles(node *FileNode, output string, depth int) string {
	if node != fb.rootNode {
		output += renderEntry(node, depth)
	}
	for _, dir := range node.subDirectories {
		output = fb.stringifyDirsAndFiles(dir, output, depth+1)
	}
	for _, file := range node.subFiles {
		output = fb.stringifyDirsAndFiles(file, output, depth+1)
	}

	return output
}

func renderEntry(entry *FileNode, depth int) string {
	var output string
	for i := 1; i < depth; i++ {
		output += " "
	}
	output += entry.Name
	if entry.IsDir {
		output += "/"
	}
	output += fmt.Sprint(LF)
	return output
}
