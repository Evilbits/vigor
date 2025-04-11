# VI(gor)

VIgor is a fun side project of implementing a terminal based IDE written in Go based on Vim behaviour. This is my first project in Go so it therefore also acts as my entrypoint into coding in Go.


As this is a pet project I specifically wanted to implement the UI with a low-level library instead of using libraries such as [tview](https://github.com/rivo/tview). This project exists only as a way to try out Go through a fun side project.

## Usage
To run from source simply execute

`go run . -f yourfile`

## Concepts
1. `Editor` is the top level component. It owns and manages all IDE features such as reading files, dealing with key inputs, etc.
2. `Screen` wraps a `Grid` (see below) and is in charge of the main event loop which renders components.
3. `Grid` wraps all renderable components and automatically sizes itself according to the CLI termsize.
4. `Drawable` components can be rendered with output. There are multiple types of `Drawable` components with different behaviour.

## Drawable
* `TextArea`: For rendering and interacting with text.
* `StatusBar`: For displaying additional metadata in the main IDE view.
* `Cmd`: Used to run commands. Behaves similar to VIM command-line mode.
* `FileBrowser`: Renders and interacts with the underlying filesystem.

## TODO
There are lots of features that still need to be implemented. This is a non-exhaustive list:
* Linenumbers in `TextArea`
* Buffer file reads
* When reading a file store text as a [Rope](https://en.wikipedia.org/wiki/Rope_(data_structure))
* Undo/redo
* x axis scroll within file
* LSP integration and syntax highlighting
* Reading a config when initiating the `Editor`

## Implemented behaviour
* Smart cursor behaviour.
    * Keeps track of x position when moving between lines of different length.
* Navigating a full file with y axis scroll (currently x axis is not supported).
* Inserting and removing characters. 
* Full browsing of a repository with unlimited file nesting.

## Implemented commands
# Visual mode
* `h`, `j`, `k`, `l`: Cursor movement
* `$`: Go to end of line
* `0`: Go to beginning of line
* `g`: Go to start of file
* `G`: Go to end of file
* `i, a`: Enter insert mode
* `:` - Enter command mode

# Command mode
* `q`, `quit`: Quit editor without saving 
* `w`, `write`: Write changes to file
* `e`, `edit`: Enter file browser
