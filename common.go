package winuserio

import "golang.org/x/sys/windows"

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

type hooksWinApi struct {
	user32              *windows.DLL
	setWindowsHookExA   *windows.Proc
	callNextHookEx      *windows.Proc
	unhookWindowsHookEx *windows.Proc
	getMessageW         *windows.Proc
}

type hookSetupResult struct {
	handle uintptr
	err    error
}
