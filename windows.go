//go:build windows

package main

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func isIndependentWindow() bool {
	// Check if started via explorer (drag-and-drop or double-click)
	// By checking the process list: if we are the only process, we own the console.
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleProcessList := kernel32.NewProc("GetConsoleProcessList")

	var processList [2]uint32
	// If the count of processes attached to this console is 1, it's an independent window
	count, _, _ := getConsoleProcessList.Call(uintptr(unsafe.Pointer(&processList)), 2)
	if count <= 1 {
		return true
	}

	// Fallback to your original mode check
	exe, _ := os.Executable()
	arg0 := strings.Trim(os.Args[0], `"`)
	if !strings.EqualFold(arg0, exe) {
		return false
	}
	
	var mode uint32
	// Use syscall explicitly to avoid ambiguity during cross-compilation
	handle := syscall.Handle(os.Stdout.Fd())
	err := syscall.GetConsoleMode(handle, &mode)
	if err != nil {
		return true 
	}
	return (mode & 0x0004) == 0
}
