// Package user32util provides helper functionality for working with
// Windows' user32 library.
//
// Many of these functions require that you first load the user32 DLL:
//	user32, err := user32util.LoadUser32DLL()
//	if err != nil {
//		// Error handling.
//	}
//
// Please refer to the Windows API documentation for more information:
// https://docs.microsoft.com/en-us/windows/win32/api/winuser/
package user32util
