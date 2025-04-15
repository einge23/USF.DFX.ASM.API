// util/usb.go
package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func FindUSBDrive() (string, error) {
	switch OnRpi {
	case false:
		return scanWindowsDrives()
	case true:
		return scanLinuxDrives()
	default:
		return "", fmt.Errorf("unsupported OS. Only Windows and Linux are supported")
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
	base := "/media/dfxp/" // Depends on the Raspberry Pi

	entries, err := os.ReadDir(base)
	if err != nil {
		return "", fmt.Errorf("failed to read /media/dfxp/: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			return filepath.Join(base, entry.Name()), nil
		}
	}
	return "", fmt.Errorf("no USB drive found on Linux")
}

func UnmountUSB() error {
	cmd := exec.Command("bash", "/home/dfxp/Desktop/AutomatedAccessControl/Repos/USF.DFX.ASM.API/scripts/unmount_all.sh")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error unmounting: %v, script output: %s", err, string(output))
	}
	return nil
}
