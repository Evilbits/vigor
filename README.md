# VI(gor)

VIgor is a fun side project of implementing a terminal based IDE written in Go based on Vim behaviour. This is my first project in Go so it therefore also acts as my entrypoint into coding in Go.


As this is a pet project I specifically wanted to implement the UI with a low-level library instead of using libraries such as [tview](https://github.com/rivo/tview)

## Usage
To run from source simply execute
`go run . -f yourfile`


## TODO
There are lots of features that still need to be implemented. This is a non-exhaustive list:
* Reading a config when initiating the `Editor`
* Inserting and removing characters
* Navigating between files
* Saving to disk
* When reading a file store text as a [Rope](https://en.wikipedia.org/wiki/Rope_(data_structure))
