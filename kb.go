package winuserio

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

// The follow code is based on work by jimmycliff obonyo:
// https://gist.github.com/obonyojimmy/52d836a1b31e2fc914d19a81bd2e0a1b

const (
	whKeyboardLl            = 13
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
	hooksWinApi *hooksWinApi
	fn          OnLowLevelKeyboardEventFunc
	hookHandle  uintptr
	done        chan error
}

func (o *LowLevelKeyboardEventListener) OnDone() chan error {
	return o.done
}

func (o *LowLevelKeyboardEventListener) Release() error {
	o.hooksWinApi.unhookWindowsHookEx.Call(o.hookHandle)

	o.hookHandle = 0

	return o.hooksWinApi.user32.Release()
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

type hooksWinApi struct {
	user32              *windows.DLL
	setWindowsHookExA   *windows.Proc
	callNextHookEx      *windows.Proc
	unhookWindowsHookEx *windows.Proc
	getMessageW         *windows.Proc
}

func NewLowLevelKeyboardListener(fn OnLowLevelKeyboardEventFunc) (*LowLevelKeyboardEventListener, error) {
	hooksWinApi, err := newHooksWinApi()
	if err != nil {
		return nil, err
	}

	var hookHandle uintptr

	ready := make(chan error)
	done := make(chan error)

	go func() {
		hookHandle, _, err = hooksWinApi.setWindowsHookExA.Call(
			uintptr(whKeyboardLl),
			uintptr(windows.NewCallback(func(nCode int, wParam uintptr, lParam uintptr) uintptr {
				if nCode == 0 {
					fn(LowLevelKeyboardEvent{
						wParam: wParam,
						lParam: lParam,
						s:      (*KbDllHookStruct)(unsafe.Pointer(lParam)),
					})
				}

				nextHookCallResult, _, _ := hooksWinApi.callNextHookEx.Call(hookHandle, uintptr(nCode), wParam, lParam)

				return nextHookCallResult
			})),
			uintptr(0),
			uintptr(0),
		)
		if hookHandle == 0 && err != nil {
			ready <- err
			return
		}

		ready <- nil

		// Needed to actually get events. Must be on same thread as hook.
		// TODO: How does this get unblocked? It's blocked forever.
		for r, _, _ := hooksWinApi.getMessageW.Call(0, 0, 0, 0); r == 0; {}

		done <- nil
	}()

	err = <-ready
	if err != nil {
		return nil, err
	}

	return &LowLevelKeyboardEventListener{
		hooksWinApi: hooksWinApi,
		hookHandle:  hookHandle,
		fn:          fn,
		done:        done,
	}, nil
}

func newHooksWinApi() (*hooksWinApi, error) {
	user32, err := windows.LoadDLL(user32DllName)
	if err != nil {
		return nil, err
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

	return &hooksWinApi{
		user32:              user32,
		setWindowsHookExA:   set,
		callNextHookEx:      call,
		unhookWindowsHookEx: unhook,
		getMessageW:         getMessageW,
	}, nil
}
