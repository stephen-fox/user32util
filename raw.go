package winuserio

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

const (
	getRawInputDeviceListName = "GetRawInputDeviceList"
)

type rawWinApi struct {
	user32                *windows.DLL
	getRawInputDeviceList *windows.Proc
}

type RawInputDeviceList struct {
	Handle uintptr
	Dword  uint32
}

func RawInputDevices() ([]RawInputDeviceList, error) {
	r, err := newRawApi()
	if err != nil {
		return []RawInputDeviceList{}, err
	}

	var numberOfDevices uint

	_, _, err = r.getRawInputDeviceList.Call(0, uintptr(unsafe.Pointer(&numberOfDevices)), unsafe.Sizeof(RawInputDeviceList{}))
	if err != nil && err.(windows.Errno) != 0 {
		return []RawInputDeviceList{}, err
	}

	devices := make([]RawInputDeviceList, numberOfDevices)

	_, _, err = r.getRawInputDeviceList.Call(uintptr(unsafe.Pointer(&devices[0])), uintptr(unsafe.Pointer(&numberOfDevices)), unsafe.Sizeof(devices[0]))
	if err != nil && err.(windows.Errno) != 0 {
		return []RawInputDeviceList{}, err
	}

	return devices, nil
}

func newRawApi() (*rawWinApi, error) {
	user32, err := windows.LoadDLL(user32DllName)
	if err != nil {
		return nil, err
	}

	get, err := user32.FindProc(getRawInputDeviceListName)
	if err != nil {
		return nil, err
	}

	return &rawWinApi{
		user32:                user32,
		getRawInputDeviceList: get,
	}, nil
}
