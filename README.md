# VI(gor)

VIgor is a fun side project of implementing a terminal based IDE written in Go based on Vim behaviour. This is my first project in Go so it therefore also acts as my entrypoint into coding in Go.


As this is a pet project I specifically wanted to implement the UI with a low-level library instead of using libraries such as [tview](https://github.com/rivo/tview).

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

## TODO
There are lots of features that still need to be implemented. This is a non-exhaustive list:
* Reading a config when initiating the `Editor`
* Linenumbers in `TextArea`
* Navigating between files
* Saving to disk
* When reading a file store text as a [Rope](https://en.wikipedia.org/wiki/Rope_(data_structure))
* Undo/redo
* x axis scroll within file

## Implemented
* Smart cursor behaviour.
    * Keeps track of x position when moving between lines of different length.
    * Vim cursor movement motions such as `$`, `0`, and going to start (`g`) and end (`G`) of file.
* Navigating a full file with y axis scroll (currently x axis is not supported).
* Inserting and removing characters. 
* Your favourite VIM commands like `:q`.
