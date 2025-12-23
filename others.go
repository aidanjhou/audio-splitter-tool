//go:build !windows

package main

import (
	"os"
	"syscall"
	"unsafe"
)

func isIndependentWindow() bool {
	var termios syscall.Termios
	// Use TIOCGWINSZ (Get Window Size) which is more universally available 
	// than TIOCGETA/TCGETS for basic TTY detection.
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, os.Stdin.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&termios)))
	
	// If it fails (err != 0), it's likely not a standard terminal (independent window)
	return err != 0
}
