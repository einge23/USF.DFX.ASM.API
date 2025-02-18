package util

import (
	"fmt"
	"strconv"
	"strings"
)

type ParsedCardData struct {
	Id int
	Username string
}

func ParseScannerString(scannerString string) (*ParsedCardData, error) {
    if (scannerString == "" || len(scannerString) >= 256) {
        return nil, fmt.Errorf("invalid scanner string format")
    }

    scannerString = strings.TrimSpace(scannerString)
    
    // Split by '^' to get the different fields
    parts := strings.Split(scannerString, "^")
    if len(parts) < 2 {
        return nil, fmt.Errorf("invalid scanner string format")
    }
    
    // Extract card number (remove the '%B' prefix)
	id := strings.TrimPrefix(parts[0], "%B")
	idNum, err := strconv.Atoi(id) 
	if err != nil {
    return nil, fmt.Errorf("failed to parse ID: %v", err)
	}    
    // Extract name and clean it up
    username := strings.TrimSpace(parts[1])
    // Replace '/' with space in name
    username = strings.ReplaceAll(username, "/", " ")
    
    return &ParsedCardData{
        Id: idNum,
        Username: username,
    }, nil
}