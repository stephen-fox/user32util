package user32util

import (
	"fmt"
	"unsafe"
)

// LowLevelMouseEvent wParam flags.
const (
	WMLButtonDown MouseButtonAction = 0x0201
	WMLButtonUp   MouseButtonAction = 0x0202
	WMMouseMove   MouseButtonAction = 0x0200
	WMMouseWheel  MouseButtonAction = 0x020A
	WMMouseHWheel MouseButtonAction = 0x020E
	WMRButtonDown MouseButtonAction = 0x0204
	WMRButtonUp   MouseButtonAction = 0x0205
)

// Other mouse related message types (unsure where they are used, but they
// appear in the 'mouseData' field documentation).
const (
	WMXButtonDown     MouseButtonAction = 0x020B
	WMXButtonUp       MouseButtonAction = 0x020C
	WMXButtonDblClk   MouseButtonAction = 0x020D
	WMNCXButtonDown   MouseButtonAction = 0x00AB
	WMNCXButtonUp     MouseButtonAction = 0x00AC
	WMNCXButtonDblClk MouseButtonAction = 0x00AD
)

// MouseButtonAction is an alias for the values contained in the
// wParam field fo LowLevelKeyboardEvent.
type MouseButtonAction uintptr

// NewLowLevelMouseListener instantiates a new mouse input listener using
// the LowLevelMouseProc Windows hook.
//
// Refer to LowLevelMouseEventListener for more information.
func NewLowLevelMouseListener(fn OnLowLevelMouseEventFunc, user32 *User32DLL) (*LowLevelMouseEventListener, error) {
	callBack := func(nCode int, wParam uintptr, lParam uintptr) {
		if nCode == 0 {
			fn(LowLevelMouseEvent{
				WParam: wParam,
				LParam: lParam,
				Struct: (*MsllHookStruct)(unsafe.Pointer(lParam)),
			})
		}
	}

	handle, threadID, done, err := setWindowsHookExW(whMouseLl, callBack, user32)
	if err != nil {
		return nil, err
	}

	return &LowLevelMouseEventListener{
		user32:     user32,
		hookHandle: handle,
		threadID:   threadID,
		fn:         fn,
		done:       done,
	}, nil
}

type OnLowLevelMouseEventFunc func(event LowLevelMouseEvent)

type LowLevelMouseEvent struct {
	WParam uintptr
	LParam uintptr
	Struct *MsllHookStruct
}

func (o LowLevelMouseEvent) MouseButtonAction() MouseButtonAction {
	return MouseButtonAction(o.WParam)
}

// From the Windows API documentation:
//
//	Contains information about a low-level mouse input event.
//
// Refer to the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msllhookstruct
type MsllHookStruct struct {
	Point       Point
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

// From the Windows API documentation:
//
//	The POINT structure defines the x- and y- coordinates of a point.
//
// Refer to the following Windows API document for more information:
// https://docs.microsoft.com/en-us/previous-versions/dd162805%28v=vs.85%29
type Point struct {
	X int32
	Y int32
}

// LowLevelMouseEventListener represents an instance of the
// LowLevelMouseProc Windows hook.
//
// From the Windows API documentation:
//
//	An application-defined or library-defined callback function
//	used with the SetWindowsHookEx function. The system calls
//	this function every time a new mouse input event is about to
//	be posted into a thread input queue.
//
// Refer to the following Windows API document for more information:
// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644986%28v=vs.85%29
type LowLevelMouseEventListener struct {
	user32     *User32DLL
	fn         OnLowLevelMouseEventFunc
	hookHandle uintptr
	threadID   uint32
	done       <-chan error
}

// OnDone returns a channel that is written to when the event listener exits.
// A non-nil error is written if an error caused the listener to exit.
func (o *LowLevelMouseEventListener) OnDone() <-chan error {
	return o.done
}

// Release releases the underlying hook handle and stops the listener from
// receiving any additional events.
func (o *LowLevelMouseEventListener) Release() error {
	o.user32.postThreadMessageW.Call(uintptr(o.threadID), wmQuit, 0, 0)

	o.user32.unhookWindowsHookEx.Call(o.hookHandle)

	o.hookHandle = 0

	return nil
}

// SetCursorPos sets the mouse cursor position.
//
// From the Windows API documentation:
//
//	Moves the cursor to the specified screen coordinates. If the new
//	coordinates are not within the screen rectangle set by the most
//	recent ClipCursor function call, the system automatically adjusts
//	the coordinates so that the cursor stays within the rectangle.
//
// Refer to the Windows API documentation for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setcursorpos
func SetCursorPos(x int32, y int32, user32 *User32DLL) (bool, error) {
	ret, _, err := user32.setCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 1 {
		return true, nil
	}

	return false, err
}

func GetCursorPos(user32 *User32DLL) (bool, error) {
	point := Point{}
	pointPointer := &point
	p := unsafe.Pointer(pointPointer)
	ret, _, err := user32.getCursorPos.Call(uintptr(p))
	if ret == 1 {
		return true, nil
	}
	fmt.Printf("The value of ret is: ", point)

	return false, err
}
