// TODO: A lot of this should be over in editor/
package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"slices"
)

type FileNode struct {
	path           string
	isDir          bool
	subFiles       []*FileNode
	subDirectories []*FileNode
}

func (f *FileNode) Name() string {
	if strings.Contains(f.path, "/") {
		splitPath := strings.Split(f.path, "/")
		return splitPath[len(splitPath)-1]
	}
	return f.path
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

func (fb *FileBrowser) GetCurrentFile() *FileNode {
	return fb.findNodeByCursorIdx()
}

func (fb *FileBrowser) loadDir(dir string) (*FileNode, error) {
	rootDir := &FileNode{
		path:           dir,
		isDir:          true,
		subFiles:       []*FileNode{},
		subDirectories: []*FileNode{},
	}
	dirItems, err := os.ReadDir(rootDir.path)
	if err != nil {
		return nil, err
	}

	for _, dirItem := range dirItems {
		node := &FileNode{
			path:           dirItem.Name(),
			isDir:          dirItem.IsDir(),
			subFiles:       []*FileNode{},
			subDirectories: []*FileNode{},
		}
		if node.isDir {
			rootDir.subDirectories = append(rootDir.subDirectories, node)
		} else {
			rootDir.subFiles = append(rootDir.subFiles, node)
		}
	}
	return rootDir, nil
}

func (fb *FileBrowser) Draw(screen *Screen) {
	fb.Box.AddText(fb.buildTextContent())
	fb.Box.Draw(screen)
	screen.RenderCursor(0, fb.cursorY, tcell.CursorStyleDefault)
}

func (fb *FileBrowser) buildTextContent() string {
	return fb.stringifyDirsAndFiles(fb.rootNode, "")
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
			return true // Found the target node
		}

		// Process directories first
		if slices.ContainsFunc(node.subDirectories, traverse) {
			return true // Target found in directory subtree
		}

		// Then process files
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

func (fb *FileBrowser) stringifyDirsAndFiles(node *FileNode, output string) string {
	if node != fb.rootNode {
		output += renderEntry(node)
	}
	for _, dir := range node.subDirectories {
		output = fb.stringifyDirsAndFiles(dir, output)
	}
	for _, file := range node.subFiles {
		output = fb.stringifyDirsAndFiles(file, output)
	}

	return output
}

func renderEntry(entry *FileNode) string {
	var output string
	output += entry.path
	if entry.isDir {
		output += "/"
	}
	output += fmt.Sprint(LF)
	return output
}
