package winuserio

import (
	"golang.org/x/sys/windows"
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
// appear in the 'mouseData' field documentation.
const (
	WMXButtonDown     MouseButtonAction = 0x020B
	WMXButtonUp       MouseButtonAction = 0x020C
	WMXButtonDblClk   MouseButtonAction = 0x020D
	WMNCXButtonDown   MouseButtonAction = 0x00AB
	WMNCXButtonUp     MouseButtonAction = 0x00AC
	WMNCXButtonDblClk MouseButtonAction = 0x00AD
)

type MouseButtonAction uintptr

type OnLowLevelMouseEventFunc func(event LowLevelMouseEvent)

type LowLevelMouseEvent struct {
	WParam uintptr
	LParam uintptr
	S      *MsllHookStruct
}

func (o LowLevelMouseEvent) MouseButtonAction() MouseButtonAction {
	return MouseButtonAction(o.WParam)
}

func NewLowLevelMouseListener(fn OnLowLevelMouseEventFunc, user32 *User32DLL) (*LowLevelMouseEventListener, error) {
	ready := make(chan hookSetupResult)
	done := make(chan error)

	go func() {
		var hookHandle uintptr
		var err error
		hookHandle, _, err = user32.setWindowsHookExA.Call(
			uintptr(whMouseLl),
			uintptr(windows.NewCallback(func(nCode int, wParam uintptr, lParam uintptr) uintptr {
				if nCode == 0 {
					fn(LowLevelMouseEvent{
						WParam: wParam,
						LParam: lParam,
						S:      (*MsllHookStruct)(unsafe.Pointer(lParam)),
					})
				}

				nextHookCallResult, _, _ := user32.callNextHookEx.Call(hookHandle, uintptr(nCode), wParam, lParam)

				return nextHookCallResult
			})),
			uintptr(0),
			uintptr(0),
		)
		if hookHandle == 0 && err != nil {
			ready <- hookSetupResult{err:err}
			return
		}

		ready <- hookSetupResult{handle: hookHandle}

		// Needed to actually get events. Must be on same thread as hook.
		// TODO: How does this get unblocked? It's blocked forever.
		for r, _, _ := user32.getMessageW.Call(0, 0, 0, 0); r == 0; {}

		done <- nil
	}()

	result := <-ready
	if result.err != nil {
		return nil, result.err
	}

	return &LowLevelMouseEventListener{
		user32:     user32,
		hookHandle: result.handle,
		fn:         fn,
		done:       done,
	}, nil
}

// See the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msllhookstruct
type MsllHookStruct struct {
	Point       Point
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type Point struct {
	X int32
	Y int32
}

// See the following Windows API document for more information:
// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/ms644986(v=vs.85)
type LowLevelMouseEventListener struct {
	user32     *User32DLL
	fn         OnLowLevelMouseEventFunc
	hookHandle uintptr
	done       chan error
}

func (o *LowLevelMouseEventListener) OnDone() chan error {
	return o.done
}

func (o *LowLevelMouseEventListener) Release() error {
	o.user32.unhookWindowsHookEx.Call(o.hookHandle)

	o.hookHandle = 0

	return nil
}
