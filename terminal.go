package main

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// detectSixelSupport checks common TERM values. It's a heuristic; we also allow
// users to force rendering later if desired.
func detectSixelSupport() bool {
	term := os.Getenv("TERM")
	support := map[string]bool{
		"xterm":          true,
		"xterm-256color": true,
		"foot":           true,
		"mlterm":         true,
		"wezterm":        true,
	}
	if support[term] {
		return true
	}
	// Also check COLORTERM hints (not definitive)
	if os.Getenv("COLORTERM") != "" {
		// not a guarantee, but many modern terminals set this
		return false
	}
	return false
}

// getTerminalSize returns columns and rows. It uses TIOCGWINSZ and returns an error
// if the ioctl fails (e.g., stdout not a TTY).
func getTerminalSize() (width, height int, err error) {
	ws := &winsize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		os.Stdout.Fd(), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 || errno != 0 {
		return 0, 0, fmt.Errorf("TIOCGWINSZ failed: %v", errno)
	}
	return int(ws.Col), int(ws.Row), nil
}
