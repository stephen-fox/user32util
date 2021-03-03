package user32util

import (
	"unsafe"
)

// LowLevelKeyboardEvent wParam flags.
const (
	WMKeyDown       KeyboardButtonAction = 256
	WMKeyUp         KeyboardButtonAction = 257
	WHSystemKeyDown KeyboardButtonAction = 260
	WMSystemKeyUp   KeyboardButtonAction = 261
)

// KeyboardButtonAction is an alias for the values contained in the
// wParam field fo LowLevelKeyboardEvent.
type KeyboardButtonAction uintptr

// NewLowLevelKeyboardListener instantiates a new keyboard input listener using
// the LowLevelKeyboardProc Windows hook.
//
// Refer to LowLevelKeyboardEventListener for more information.
func NewLowLevelKeyboardListener(fn OnLowLevelKeyboardEventFunc, user32 *User32DLL) (*LowLevelKeyboardEventListener, error) {
	callBack := func(nCode int, wParam uintptr, lParam uintptr) {
		if nCode == 0 {
			fn(LowLevelKeyboardEvent{
				WParam: wParam,
				LParam: lParam,
				Struct: (*KbdllHookStruct)(unsafe.Pointer(lParam)),
			})
		}
	}

	handle, threadID, done, err := setWindowsHookExW(whKeyboardLl, callBack, user32)
	if err != nil {
		return nil, err
	}

	return &LowLevelKeyboardEventListener{
		user32:     user32,
		hookHandle: handle,
		threadID:   threadID,
		fn:         fn,
		done:       done,
	}, nil
}

type OnLowLevelKeyboardEventFunc func(event LowLevelKeyboardEvent)

// LowLevelKeyboardEventListener represents an instance of the
// LowLevelKeyboardProc Windows hook.
//
// From the Windows API documentation:
//	An application-defined or library-defined callback function used
//	with the SetWindowsHookEx function. The system calls this function
//	every time a new keyboard input event is about to be posted into
//	a thread input queue.
//
// Refer to the following Windows API document for more information:
// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644985%28v=vs.85%29
type LowLevelKeyboardEventListener struct {
	user32     *User32DLL
	fn         OnLowLevelKeyboardEventFunc
	hookHandle uintptr
	threadID   uint32
	done       <-chan error
}

// OnDone returns a channel that is written to when the event listener exits.
// A non-nil error is written if an error caused the listener to exit.
func (o *LowLevelKeyboardEventListener) OnDone() <-chan error {
	return o.done
}

// Release releases the underlying hook handle and stops the listener from
// receiving any additional events.
func (o *LowLevelKeyboardEventListener) Release() error {
	o.user32.postThreadMessageW.Call(uintptr(o.threadID), wmQuit, 0, 0)

	o.user32.unhookWindowsHookEx.Call(o.hookHandle)

	o.hookHandle = 0

	return nil
}

// LowLevelKeyboardEvent represents a single keyboard event.
type LowLevelKeyboardEvent struct {
	WParam uintptr
	LParam uintptr
	Struct *KbdllHookStruct
}

func (o LowLevelKeyboardEvent) KeyboardButtonAction() KeyboardButtonAction {
	return KeyboardButtonAction(o.WParam)
}

// From the Windows API documentation:
//	Contains information about a low-level keyboard input event.
//
// Refer to the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-kbdllhookstruct
type KbdllHookStruct struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

func (o KbdllHookStruct) VirtualKeyCode() byte {
	return byte(o.VkCode)
}
