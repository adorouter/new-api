package service

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdoDiscordAccountOldEnough(t *testing.T) {
	now := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	idCreatedAt := func(at time.Time) string {
		return strconv.FormatUint(uint64(at.UnixMilli()-adoDiscordEpochMillis)<<22, 10)
	}

	tests := []struct {
		name string
		id   string
		want bool
	}{
		{"five years old", idCreatedAt(now.AddDate(-5, 0, 0)), true},
		{"31 days old", idCreatedAt(now.Add(-31 * 24 * time.Hour)), true},
		{"29 days old", idCreatedAt(now.Add(-29 * 24 * time.Hour)), false},
		{"made today", idCreatedAt(now.Add(-time.Hour)), false},
		{"not a number", "not-a-snowflake", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, adoDiscordAccountOldEnough(tt.id, now))
		})
	}
}
