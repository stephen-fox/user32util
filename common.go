package user32util

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	wmQuit = 0x0012
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

	set, err := user32.FindProc(setWindowsHookExAName)
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

	return &User32DLL{
		user32:              user32,
		setWindowsHookExA:   set,
		callNextHookEx:      call,
		unhookWindowsHookEx: unhook,
		getMessageW:         getMessageW,
		sendInput:           sendInput,
		postThreadMessageW:  postThreadMessageW,
	}, nil
}

// User32DLL represents the user32 DLL, mapping several of its procedures to
// this struct's fields.
type User32DLL struct {
	user32              *windows.DLL
	setWindowsHookExA   *windows.Proc
	callNextHookEx      *windows.Proc
	unhookWindowsHookEx *windows.Proc
	getMessageW         *windows.Proc
	sendInput           *windows.Proc
	postThreadMessageW  *windows.Proc
}

type hookSetupResult struct {
	handle uintptr
	tid    uint32
	err    error
}

// From the Windows API documentation:
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
