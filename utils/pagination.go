package utils

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DecodeCursor decodes the cursor to get timestamp and transaction ID
func DecodeCursor(cursor string) (time.Time, uint, error) {
	decoded, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, 0, err
	}

	parts := strings.Split(string(decoded), "_")
	if len(parts) != 2 {
		return time.Time{}, 0, fmt.Errorf("invalid cursor format")
	}

	timestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, 0, err
	}
	txnID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return time.Time{}, 0, err
	}

	return time.Unix(0, timestamp*int64(time.Millisecond)), uint(txnID), nil
}

// EncodeCursor encodes the timestamp and transaction ID into a cursor
func EncodeCursor(timestamp time.Time, transactionID uint) string {
	cursor := fmt.Sprintf("%d_%d", timestamp.UnixMilli(), transactionID)
	return base64.StdEncoding.EncodeToString([]byte(cursor))
}
