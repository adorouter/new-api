package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAdoDiscordAccess(t *testing.T) {
	const guildID = "1561373731224555630"

	tests := []struct {
		name         string
		guildID      string
		verified     bool
		email        string
		memberStatus int
		wantMessage  string
		wantCalled   bool
	}{
		{name: "verified member signs in", guildID: guildID, verified: true, email: "a@example.com", memberStatus: http.StatusOK, wantCalled: true},
		{name: "non-member is told to join", guildID: guildID, verified: true, email: "a@example.com", memberStatus: http.StatusNotFound, wantMessage: adoDiscordJoinMessage, wantCalled: true},
		{name: "unverified email is refused before the member check", guildID: guildID, verified: false, email: "a@example.com", memberStatus: http.StatusOK, wantMessage: adoDiscordEmailMessage},
		{name: "missing email is refused", guildID: guildID, verified: true, email: "", memberStatus: http.StatusOK, wantMessage: adoDiscordEmailMessage},
		{name: "missing scope fails closed", guildID: guildID, verified: true, email: "a@example.com", memberStatus: http.StatusForbidden, wantMessage: adoDiscordCheckMessage, wantCalled: true},
		{name: "rate limit fails closed", guildID: guildID, verified: true, email: "a@example.com", memberStatus: http.StatusTooManyRequests, wantMessage: adoDiscordCheckMessage, wantCalled: true},
		{name: "outage fails closed", guildID: guildID, verified: true, email: "a@example.com", memberStatus: http.StatusBadGateway, wantMessage: adoDiscordCheckMessage, wantCalled: true},
		{name: "unset server ID fails closed", guildID: "", verified: true, email: "a@example.com", memberStatus: http.StatusOK, wantMessage: adoDiscordConfigMessage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				assert.Equal(t, "/users/@me/guilds/"+guildID+"/member", r.URL.Path)
				assert.Equal(t, "Bearer user-token", r.Header.Get("Authorization"))
				w.WriteHeader(tt.memberStatus)
			}))
			defer server.Close()

			oldBase := adoDiscordAPIBase
			adoDiscordAPIBase = server.URL
			t.Cleanup(func() { adoDiscordAPIBase = oldBase })
			t.Setenv("DISCORD_GUILD_ID", tt.guildID)

			err := checkAdoDiscordAccess(context.Background(), "user-token", tt.verified, tt.email)

			assert.Equal(t, tt.wantCalled, called)
			if tt.wantMessage == "" {
				require.NoError(t, err)
				return
			}
			var denied *AccessDeniedError
			require.ErrorAs(t, err, &denied)
			assert.Equal(t, tt.wantMessage, denied.Message)
		})
	}
}

func TestCheckAdoDiscordAccessFailsClosedWhenDiscordIsUnreachable(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	unreachable := server.URL
	server.Close()

	oldBase := adoDiscordAPIBase
	adoDiscordAPIBase = unreachable
	t.Cleanup(func() { adoDiscordAPIBase = oldBase })
	t.Setenv("DISCORD_GUILD_ID", "1561373731224555630")

	err := checkAdoDiscordAccess(context.Background(), "user-token", true, "a@example.com")

	var denied *AccessDeniedError
	require.ErrorAs(t, err, &denied)
	assert.Equal(t, adoDiscordCheckMessage, denied.Message)
}
