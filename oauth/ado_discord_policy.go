package oauth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/logger"
)

// adoDiscordAPIBase is a variable so tests can point it at a local server.
var adoDiscordAPIBase = "https://discord.com/api/v10"

const (
	adoDiscordEmailMessage   = "Verify the email address on your Discord account, then sign in again."
	adoDiscordJoinMessage    = "Join the adorouter Discord server, then sign in again."
	adoDiscordCheckMessage   = "We couldn't check your Discord server membership. Please try again in a moment."
	adoDiscordConfigMessage  = "Sign-in isn't available right now. Please try again later."
	adoDiscordMemberCheckFmt = "%s/users/@me/guilds/%s/member"
)

// checkAdoDiscordAccess enforces adorouter's sign-in rules on every Discord
// sign-in: a verified email and membership in the adorouter server. It fails
// closed, so a missing server ID or a Discord outage never lets anyone through.
func checkAdoDiscordAccess(ctx context.Context, accessToken string, verified bool, email string) error {
	guildID := strings.TrimSpace(os.Getenv("DISCORD_GUILD_ID"))
	if guildID == "" {
		logger.LogError(ctx, common.LogText("[OAuth-Discord] DISCORD_GUILD_ID is not set, refusing sign-in"))
		return &AccessDeniedError{Message: adoDiscordConfigMessage}
	}
	if !verified || strings.TrimSpace(email) == "" {
		return &AccessDeniedError{Message: adoDiscordEmailMessage}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(adoDiscordMemberCheckFmt, adoDiscordAPIBase, guildID), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		logger.LogError(ctx, common.LogText("[OAuth-Discord] member check error: %s", err.Error()))
		return &AccessDeniedError{Message: adoDiscordCheckMessage}
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return &AccessDeniedError{Message: adoDiscordJoinMessage}
	default:
		logger.LogError(ctx, common.LogText("[OAuth-Discord] member check failed: status=%d", res.StatusCode))
		return &AccessDeniedError{Message: adoDiscordCheckMessage}
	}
}
