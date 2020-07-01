# user32util

[![GoDoc][godoc-badge]][godoc]

[godoc-badge]: https://godoc.org/github.com/stephen-fox/user32util?status.svg
[godoc]: https://godoc.org/github.com/stephen-fox/user32util

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

- [readkeyboard](examples/readkeyboard/main.go) - Reads keyboard presses and
prints them to stderr
- [readmouse](examples/readmouse/main.go) - Reads mouse inputs and prints them
to stderr
- [sendinput](examples/sendinput/main.go) - Sends keyboard or mouse inputs
to Windows
