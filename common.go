package winuserio

import (
	"golang.org/x/sys/windows"
)

func loadUser32DLL() (*user32DLL, error) {
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

	return &user32DLL{
		user32:              user32,
		setWindowsHookExA:   set,
		callNextHookEx:      call,
		unhookWindowsHookEx: unhook,
		getMessageW:         getMessageW,
	}, nil
}

type user32DLL struct {
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
