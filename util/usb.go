package util

import (
	"database/sql"
	"encoding/csv"
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

// given the name of a db table and the name of a proposed csv file, create a .csv file
// that represents that table in the db, with the name provided.
func ExportTableToCSV(tableName string, outputCSV string) error {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %v", err)
	}

	file, err := os.Create(outputCSV)
	if err != nil {
		return fmt.Errorf("failed to create csv file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(columns); err != nil {
		return fmt.Errorf("failed to write header: %v", err)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}

		record := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				record[i] = ""
			} else {
				record[i] = fmt.Sprintf("%v", val)
			}
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write row: %v", err)
		}
	}

	return nil
}

// given a string to represent a destination, copy a file to that destination.
func MoveFile(source string, destination string) error {
	cmd := exec.Command("cp", source, destination)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error moving db file from %s to %s: %v, command output: %s", source, destination, err, string(output))
	}
	return nil
}
