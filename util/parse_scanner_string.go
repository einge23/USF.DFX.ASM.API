package util

import (
	"fmt"
	"strconv"
	"strings"
)

type ParsedCardData struct {
	Id       int
	Username string
}

func ParseScannerString(scannerString string) (*ParsedCardData, error) {
	if scannerString == "" || len(scannerString) >= 256 {
		return nil, fmt.Errorf("invalid scanner string format")
	}

	scannerString = strings.TrimSpace(scannerString)

	// Split by '^' to get the different fields
	parts := strings.Split(scannerString, "^")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid scanner string format")
	}

	id := strings.TrimPrefix(parts[0], "%B")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %v", err)
	}
	username := strings.TrimSpace(parts[1])
	username = strings.ReplaceAll(username, "/", " ")

	nameParts := strings.Fields(username)
	if len(nameParts) >= 2 {

		firstName := nameParts[1]
		lastName := nameParts[0]
		middleNames := nameParts[2:]

		allParts := append([]string{firstName}, middleNames...)
		allParts = append(allParts, lastName)
		username = strings.Join(allParts, " ")
	}

	return &ParsedCardData{
		Id:       idNum,
		Username: username,
	}, nil
}
