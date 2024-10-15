package detect

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                       = windows.NewLazySystemDLL("user32.dll")
	procSendMessage              = user32.NewProc("SendMessageW")
	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
)

const (
	WM_CLOSE = 0x0010
)

func GetProcessPathByName(processName string) (string, error) {
	pid, err := getPIDByName(processName)
	if err != nil {
		return "", err
	}

	handle, err := openProcess(pid)
	if err != nil {
		return "", err
	}

	defer windows.CloseHandle(handle)

	exePath, err := getProcessExePath(handle)
	if err != nil {
		return "", err
	}

	return exePath, nil
}

func CloseProcessWindow(processName string) error {
	hwnd, err := findWindow(processName)
	if err != nil {
		return err
	}

	_, err = sendMessage(hwnd, WM_CLOSE, 0, 0)
	if err != nil && err != syscall.EINVAL {
		return err
	}

	return nil
}

func getPIDByName(processName string) (uint32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(
		windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, err
	}

	defer windows.CloseHandle(snapshot)

	var pe32 windows.ProcessEntry32
	pe32.Size = uint32(unsafe.Sizeof(pe32))

	err = windows.Process32First(snapshot, &pe32)
	if err != nil {
		return 0, err
	}

	for {
		name := windows.UTF16ToString(pe32.ExeFile[:])
		if name == processName {
			return pe32.ProcessID, nil
		}

		err = windows.Process32Next(snapshot, &pe32)
		if err != nil {
			break
		}
	}

	return 0, fmt.Errorf("not found: %s", processName)
}

func openProcess(pid uint32) (windows.Handle, error) {
	handle, err := windows.OpenProcess(
		windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil {
		return 0, err
	}

	return handle, nil
}

func getProcessExePath(handle windows.Handle) (string, error) {
	var buf [windows.MAX_PATH]uint16
	size := uint32(len(buf))

	err := windows.QueryFullProcessImageName(handle, 0, &buf[0], &size)
	if err != nil {
		return "", err
	}

	return syscall.UTF16ToString(buf[:]), nil
}

func findWindow(processName string) (windows.Handle, error) {
	pid, err := getPIDByName(processName)
	if err != nil {
		return 0, err
	}

	var hwnd windows.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		var procID uint32
		procGetWindowThreadProcessId.Call(uintptr(h), uintptr(unsafe.Pointer(&procID)))
		if procID == pid {
			hwnd = windows.Handle(h)
			return 0
		}
		return 1
	})

	procEnumWindows.Call(cb, 0)
	if hwnd == 0 {
		return 0, fmt.Errorf("no window found for process %s", processName)
	}
	return hwnd, nil
}

func sendMessage(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) (lResult uintptr, err error) {
	r0, _, e1 := syscall.Syscall6(procSendMessage.Addr(), 4, uintptr(hwnd), uintptr(msg), wParam, lParam, 0, 0)
	if r0 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	lResult = r0
	return
}
