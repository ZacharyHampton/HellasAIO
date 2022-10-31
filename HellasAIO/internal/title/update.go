package title

import (
	"fmt"
	"syscall"
	"unsafe"
)

func setConsoleTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}
	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}

func updateTitle() {
	_, _ = setConsoleTitle(fmt.Sprintf("HellasAIO ｜ Carts: %d ｜ Checkouts: %d ｜ Failures: %d", carts, checkouts, failures))
}
