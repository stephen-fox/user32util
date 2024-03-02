package user32util

import (
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Various WM codes.
const (
	wmQuit = 0x0012
)

const (
	whKeyboardLl            = 13
	whMouseLl               = 14
	user32DllName           = "user32.dll"
	setWindowsHookExWName   = "SetWindowsHookExW"
	callNextHookExName      = "CallNextHookEx"
	unhookWindowsHookExName = "UnhookWindowsHookEx"
	getMessageWName         = "GetMessageW"
	sendInputName           = "SendInput"
	postThreadMessageWName  = "PostThreadMessageW"
	setCursorPosName        = "SetCursorPos"
	getCursorPosName        = "GetCursorPos"
)

// LoadUser32DLL loads the user32 DLL into memory.
func LoadUser32DLL() (*User32DLL, error) {
	// TODO: Hack to avoid using unsafe 'windows.LoadDLL()' while
	//  retaining full control over when a DLL is loaded.
	temp := windows.LazyDLL{
		Name:   user32DllName,
		System: true,
	}
	err := temp.Load()
	if err != nil {
		return nil, err
	}

	user32 := &windows.DLL{
		Name:   temp.Name,
		Handle: windows.Handle(temp.Handle()),
	}

	setWindowsHookExW, err := user32.FindProc(setWindowsHookExWName)
	if err != nil {
		return nil, err
	}

	call, err := user32.FindProc(callNextHookExName)
	if err != nil {
		return nil, err
	}

	unhook, err := user32.FindProc(unhookWindowsHookExName)
	if err != nil {
		return nil, err
	}

	getMessageW, err := user32.FindProc(getMessageWName)
	if err != nil {
		return nil, err
	}

	sendInput, err := user32.FindProc(sendInputName)
	if err != nil {
		return nil, err
	}

	postThreadMessageW, err := user32.FindProc(postThreadMessageWName)
	if err != nil {
		return nil, err
	}

	setCursorPos, err := user32.FindProc(setCursorPosName)
	if err != nil {
		return nil, err
	}

	getCursorPos, err := user32.FindProc(getCursorPosName)
	if err != nil {
		return nil, err
	}

	return &User32DLL{
		user32:              user32,
		setWindowsHookExW:   setWindowsHookExW,
		callNextHookEx:      call,
		unhookWindowsHookEx: unhook,
		getMessageW:         getMessageW,
		sendInput:           sendInput,
		postThreadMessageW:  postThreadMessageW,
		setCursorPos:        setCursorPos,
		getCursorPos:        getCursorPos,
	}, nil
}

// User32DLL represents the user32 DLL, mapping several of its procedures to
// this struct's fields.
type User32DLL struct {
	user32              *windows.DLL
	setWindowsHookExW   *windows.Proc
	callNextHookEx      *windows.Proc
	unhookWindowsHookEx *windows.Proc
	getMessageW         *windows.Proc
	sendInput           *windows.Proc
	postThreadMessageW  *windows.Proc
	setCursorPos        *windows.Proc
	getCursorPos        *windows.Proc
}

// Release releases the underlying DLL.
func (o *User32DLL) Release() error {
	return o.user32.Release()
}

// onHookCalledFunc defines what happens when a Windows hook created using
// "SetWindowsHookEx*()" is called.
type onHookCalledFunc func(nCode int, wParam uintptr, lParam uintptr)

// setWindowsHookExW wraps the 'SetWindowsHookExW()' system call, creating
// a new Windows hook for the given hook ID and callback. On success,
// it returns a handle to the hook, the ID of thread associated with the hook,
// and a channel that is written to when the hook exits.
//
// From the Windows API documentation:
//
//	Installs an application-defined hook procedure into a hook chain.
//	You would install a hook procedure to monitor the system for certain
//	types of events. These events are associated either with a specific
//	thread or with all threads in the same desktop as the calling thread.
//
// Refer to the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setwindowshookexw
func setWindowsHookExW(hookID int, callBack onHookCalledFunc, user32 *User32DLL) (uintptr, uint32, <-chan error, error) {
	ready := make(chan hookSetupResult)
	done := make(chan error, 1)

	go func() {
		runtime.LockOSThread()

		var hookHandle uintptr
		var err error
		hookHandle, _, err = user32.setWindowsHookExW.Call(
			uintptr(hookID),
			windows.NewCallback(func(nCode int, wParam uintptr, lParam uintptr) uintptr {
				callBack(nCode, wParam, lParam)

				nextHookCallResult, _, _ := user32.callNextHookEx.Call(hookHandle, uintptr(nCode), wParam, lParam)

				return nextHookCallResult
			}),
			0,
			0,
		)
		if hookHandle == 0 {
			ready <- hookSetupResult{err: err}
			return
		}

		ready <- hookSetupResult{
			handle: hookHandle,
			tid:    windows.GetCurrentThreadId(),
		}

		// Needed to actually get events. Must be on same thread as hook.
		var msg Msg
		for r, _, _ := user32.getMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0); r != 0; {
		}

		done <- nil
	}()

	result := <-ready
	if result.err != nil {
		return 0, 0, nil, result.err
	}

	return result.handle, result.tid, done, nil
}

type hookSetupResult struct {
	handle uintptr
	tid    uint32
	err    error
}

// From the Windows API documentation:
//
//	Contains message information from a thread's message queue.
//
// Refer to the following Windows API document for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/ns-winuser-msg
type Msg struct {
	Hwnd     unsafe.Pointer
	Message  uint
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}
