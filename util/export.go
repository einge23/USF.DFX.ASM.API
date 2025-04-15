package util

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

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

// given a string to represent a path, copy the database to that path.
func ExportDbFile(path string) error {
	cmd := exec.Command("cp", "/home/dfxp/Desktop/AutomatedAccessControl/Repos/USF.DFX.ASM.API/test.db", path)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error copying db file to USB drive: %v, command output: %s", err, string(output))
	}
	return nil
}
