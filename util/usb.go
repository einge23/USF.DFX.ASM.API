// util/usb.go
package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func FindUSBDrive() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return scanWindowsDrives()
	case "linux":
		return scanLinuxDrives()
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

func scanWindowsDrives() (string, error) {
	drives := []string{"E:", "D:", "F:", "G:", "H:"} // Common USB drive letters
	for _, drive := range drives {
		if info, err := os.Stat(drive + "\\"); err == nil && info.IsDir() {
			return drive + "\\", nil
		}
	}
	return "", fmt.Errorf("no USB drive found on Windows")
}

func scanLinuxDrives() (string, error) {
	base := "/media/pi" // Depends on the Raspberry Pi

	entries, err := os.ReadDir(base)
	if err != nil {
		return "", fmt.Errorf("failed to read /media: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			return filepath.Join(base, entry.Name()), nil
		}
	}
	return "", fmt.Errorf("no USB drive found on Linux")
}
