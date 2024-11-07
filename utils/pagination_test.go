package utils_test

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wanliqun/go-wallet-app/utils"
)

func TestEncodeDecodeCursor(t *testing.T) {
	timestamp := time.Now()
	transactionID := uint(12345)

	cursor := utils.EncodeCursor(timestamp, transactionID)
	assert.NotEmpty(t, cursor, "encoded cursor should not be empty")

	decodedTimestamp, decodedTransactionID, err := utils.DecodeCursor(cursor)
	assert.NoError(t, err, "decode cursor should not return an error")
	assert.Equal(t, timestamp.UnixMilli(), decodedTimestamp.UnixMilli(), "decoded timestamp should match original")
	assert.Equal(t, transactionID, decodedTransactionID, "decoded transaction ID should match original")
}

func TestDecodeCursorInvalidFormat(t *testing.T) {
	invalidCursor := base64.StdEncoding.EncodeToString([]byte("abcd1234"))

	_, _, err := utils.DecodeCursor(invalidCursor)
	assert.Error(t, err, "decoding an invalid cursor should return an error")
	assert.Contains(t, err.Error(), "invalid cursor format", "error message should indicate invalid format")
}

func TestDecodeCursorBase64Error(t *testing.T) {
	invalidBase64Cursor := "!!!not_base64_encoded"

	_, _, err := utils.DecodeCursor(invalidBase64Cursor)
	assert.Error(t, err, "decoding a non-Base64 cursor should return an error")
}
