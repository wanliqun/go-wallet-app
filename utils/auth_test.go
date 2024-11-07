package utils

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestExtractBearerToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		authorization string
		expectedToken string
		expectedError error
	}{
		{
			name:          "Valid Bearer Token",
			authorization: "Bearer abcdef12345",
			expectedToken: "abcdef12345",
			expectedError: nil,
		},
		{
			name:          "Empty Authorization Header",
			authorization: "",
			expectedToken: "",
			expectedError: ErrInvalidAuthorizationHeader,
		},
		{
			name:          "Invalid Prefix",
			authorization: "Token abcdef12345",
			expectedToken: "",
			expectedError: ErrInvalidAuthorizationHeader,
		},
		{
			name:          "Bearer Without Token",
			authorization: "Bearer ",
			expectedToken: "",
			expectedError: ErrInvalidAuthorizationHeader,
		},
		{
			name:          "Bearer with Extra Spaces",
			authorization: "   Bearer abcdef12345  ",
			expectedToken: "abcdef12345",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(nil)
			c.Request = &http.Request{Header: make(http.Header)}
			c.Request.Header.Set("Authorization", test.authorization)

			token, err := ExtractBearerToken(c)

			assert.Equal(t, test.expectedToken, token)
			if test.expectedError != nil {
				assert.True(t, errors.Is(err, test.expectedError))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
