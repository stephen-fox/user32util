package winuserio

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

// The follow code is based on work by jimmycliff obonyo:
// https://gist.github.com/obonyojimmy/52d836a1b31e2fc914d19a81bd2e0a1b

const (
	whKeyboardLl            = 13
	whMouseLl               = 14
	user32DllName           = "user32.dll"
	setWindowsHookExAName   = "SetWindowsHookExA"
	callNextHookExName      = "CallNextHookEx"
	unhookWindowsHookExName = "UnhookWindowsHookEx"
	getMessageWName         = "GetMessageW"
)

type KeyboardButtonAction uintptr

const (
	WMKeyDown       KeyboardButtonAction = 256
	WMKeyUp         KeyboardButtonAction = 257
	WHSystemKeyDown KeyboardButtonAction = 260
	WMSystemKeyUp   KeyboardButtonAction = 261
)

type OnLowLevelKeyboardEventFunc func(event LowLevelKeyboardEvent)

type LowLevelKeyboardEventListener struct {
	user32     *User32DLL
	fn         OnLowLevelKeyboardEventFunc
	hookHandle uintptr
	done       chan error
}

func (o *LowLevelKeyboardEventListener) OnDone() chan error {
	return o.done
}

func (o *LowLevelKeyboardEventListener) Release() error {
	o.user32.unhookWindowsHookEx.Call(o.hookHandle)

	o.hookHandle = 0

	return nil
}

type LowLevelKeyboardEvent struct {
	wParam uintptr
	lParam uintptr
	s      *KbDllHookStruct
}

func (o LowLevelKeyboardEvent) KeyboardButtonAction() KeyboardButtonAction {
	return KeyboardButtonAction(o.wParam)
}

func (o LowLevelKeyboardEvent) HookStruct() *KbDllHookStruct {
	return o.s
}

type KbDllHookStruct struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

func (o KbDllHookStruct) VirtualKeyCode() byte {
	return byte(o.VkCode)
}

func NewLowLevelKeyboardListener(fn OnLowLevelKeyboardEventFunc, user32 *User32DLL) (*LowLevelKeyboardEventListener, error) {
	ready := make(chan hookSetupResult)
	done := make(chan error)

	go func() {
		var hookHandle uintptr
		var err error
		hookHandle, _, err = user32.setWindowsHookExA.Call(
			uintptr(whKeyboardLl),
			uintptr(windows.NewCallback(func(nCode int, wParam uintptr, lParam uintptr) uintptr {
				if nCode == 0 {
					fn(LowLevelKeyboardEvent{
						wParam: wParam,
						lParam: lParam,
						s:      (*KbDllHookStruct)(unsafe.Pointer(lParam)),
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

		ready <- hookSetupResult{handle:hookHandle}

		// Needed to actually get events. Must be on same thread as hook.
		// TODO: How does this get unblocked? It's blocked forever.
		for r, _, _ := user32.getMessageW.Call(0, 0, 0, 0); r == 0; {}

		done <- nil
	}()

	result := <-ready
	if result.err != nil {
		return nil, result.err
	}

	return &LowLevelKeyboardEventListener{
		user32:     user32,
		hookHandle: result.handle,
		fn:         fn,
		done:       done,
	}, nil
}
