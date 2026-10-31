package oauth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckAdoDiscordAccess(t *testing.T) {
	tests := []struct {
		name     string
		verified bool
		email    string
		allowed  bool
	}{
		{name: "verified email signs in", verified: true, email: "a@example.com", allowed: true},
		{name: "unverified email is refused", verified: false, email: "a@example.com", allowed: false},
		{name: "missing email is refused", verified: true, email: " ", allowed: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkAdoDiscordAccess(tt.verified, tt.email)
			if tt.allowed {
				assert.NoError(t, err)
				return
			}
			var denied *AccessDeniedError
			assert.True(t, errors.As(err, &denied))
			assert.Equal(t, adoDiscordEmailMessage, denied.Message)
		})
	}
}
