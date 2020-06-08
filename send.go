package winuserio

import (
	"fmt"
	"unsafe"
)

const (
	InputMouse = iota
	InputKeyboard
	InputHardware
)

// Various MouseInput dwFlags.
//
// See the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-mouseinput
const (
	MouseEventFAbsolute       uint32 = 0x8000
	MouseEventFHWheel         uint32 = 0x01000
	MouseEventFMove           uint32 = 0x0001
	MouseEventFMoveNoCoalesce uint32 = 0x2000
	MouseEventFLeftDown       uint32 = 0x0002
	MouseEventFLeftUp         uint32 = 0x0004
	MouseEventFRightDown      uint32 = 0x0008
	MouseEventFRightUp        uint32 = 0x0010
	MouseEventFMiddleDown     uint32 = 0x0020
	MouseEventFMiddleUp       uint32 = 0x0040
	MouseEventFVirtualDesk    uint32 = 0x4000
	MouseEventFWheel          uint32 = 0x0800
	MouseEventFXDown          uint32 = 0x0080
	MouseEventFXUp            uint32 = 0x0100
)

// Various KeybdInput dwFlags.
//
// See the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-keybdinput
const (
	KeyEventFExtendedKey uint32 = 0x0001
	KeyEventFKeyUp       uint32 = 0x0002
	KeyEventFScanCode    uint32 = 0x0008
	KeyEventFUnicode     uint32 = 0x0004
)

// See the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-mouseinput
type MouseInput struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

func SendMouseInput(input MouseInput, user32 *User32DLL) error {
	// Apparently, no byte padding is needed.
	s := struct {
		Type uint32
		Val  MouseInput
	}{
		Type: InputMouse,
		Val:  input,
	}
	return SendInput(unsafe.Pointer(&s), unsafe.Sizeof(s), user32)
}

// See the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-keybdinput
type KeybdInput struct {
	WVK         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

func SendKeydbInput(input KeybdInput, user32 *User32DLL) error {
	// uint64's worth of padding needed to make Windows happy.
	// This is so we can omit the other structs, thereby making go
	// AND Windows happy.
	s := struct {
		Type uint32
		Val  KeybdInput
		Padd uint64
	}{
		Type: InputKeyboard,
		Val:  input,
		Padd: 0,
	}
	return SendInput(unsafe.Pointer(&s), unsafe.Sizeof(s), user32)
}

// See the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-hardwareinput
type HardwareInput struct {
	UMsg    uint32
	WParamL uint16
	WParamH uint16
}

// No idea if this works. Untested.
func SendHardwareInput(input HardwareInput, user32 *User32DLL) error {
	s := struct {
		Type uint32
		Val  HardwareInput
		Padd uint64
	}{
		Type: InputHardware,
		Val:  input,
		Padd: 0,
	}
	return SendInput(unsafe.Pointer(&s), unsafe.Sizeof(s), user32)
}

// Hacky implementation of SendInput that works around lack of union support.
//
// https://github.com/JamesHovious/w32/blob/master/user32.go works around this
// by using cgo. I have no desire to made cgo a dependency of the project.
//
// See the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-input
func SendInput(unsafePointerToVal unsafe.Pointer, inputStructSizeBytes uintptr, user32 *User32DLL) error {
	numSent, _, err := user32.sendInput.Call(
		uintptr(1),
		uintptr(unsafePointerToVal),
		uintptr(inputStructSizeBytes))
	if numSent == 1 {
		return nil
	} else if err != nil {
		return err
	}

	return fmt.Errorf("failed to send input, unknown errror")
}
