package oauth

import "strings"

const adoDiscordEmailMessage = "Verify the email address on your Discord account, then sign in again."

// checkAdoDiscordAccess is adorouter's rule for Discord sign-in: a verified email,
// because the one-time signup credit relies on Discord accounts being real.
func checkAdoDiscordAccess(verified bool, email string) error {
	if !verified || strings.TrimSpace(email) == "" {
		return &AccessDeniedError{Message: adoDiscordEmailMessage}
	}
	return nil
}
