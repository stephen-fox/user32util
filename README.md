# user32util

[![GoDoc][godoc-badge]][godoc]

[godoc-badge]: https://pkg.go.dev/badge/github.com/stephen-fox/user32util
[godoc]: https://pkg.go.dev/github.com/stephen-fox/user32util

Package user32util provides helper functionality for working with Windows'
user32 library.

## APIs
The library offers several helper functions for working with user32.

Many of these functions require that you first load the user32 DLL:
```go
user32, err := user32util.LoadUser32DLL()
if err != nil {
	// Error handling.
}
```

#### Input listeners

- `NewLowLevelMouseListener()` - Starts a listener that reports on mouse input
- `NewLowLevelKeyboardListener()` - Starts a listener that reports on
keyboard input

#### Send input

- `SendKeydbInput()` - Sends a single keyboard input
- `SendMouseInput()` - Sends a single mouse input
- `SendInput()` - Send input implements the `SendInput()` Windows system call
- `SendHardwareInput()` - Sends a single hardware input

## Examples
The following examples can be found in the [examples/ directory](examples/):

- [moveandclickmouse](examples/moveandclickmouse/main.go) - Moves the mouse
and then left clicks on the new position. Takes inputs as command line
arguments in `x,y` format. E.g., `example 1221,244 460,892`. Coordinates
can be printed by running: `example print`
- [readkeyboard](examples/readkeyboard/main.go) - Reads keyboard presses and
prints them to stderr
- [readmouse](examples/readmouse/main.go) - Reads mouse inputs and prints them
to stderr
- [sendinput](examples/sendinput/main.go) - Sends keyboard or mouse inputs
to Windows

## Special thanks
This library is influenced by jimmycliff obonyo's work in this GitHub gist:
https://gist.github.com/obonyojimmy/52d836a1b31e2fc914d19a81bd2e0a1b

Thank you for documenting your work, jimmycliff.
