package service

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
)

const (
	adoDiscordEpochMillis      = 1420070400000
	adoSignupGrantMinAge       = 30 * 24 * time.Hour
	adoSignupGrantDefaultUSD   = 0.5
	adoSignupGrantDefaultLimit = 20
	adoAlertCooldownSeconds    = 6 * 60 * 60
)

var adoLastCapAlertUnix atomic.Int64

// adoDiscordAccountOldEnough reads the account's creation time from its Discord
// ID, so a freshly made alt account gets no credit. An unreadable ID gets none.
func adoDiscordAccountOldEnough(discordID string, now time.Time) bool {
	id, err := strconv.ParseUint(strings.TrimSpace(discordID), 10, 64)
	if err != nil {
		return false
	}
	createdAt := time.UnixMilli(int64(id>>22) + adoDiscordEpochMillis)
	return now.Sub(createdAt) >= adoSignupGrantMinAge
}

func adoEnvFloat(name string, fallback float64) float64 {
	value, err := strconv.ParseFloat(strings.TrimSpace(os.Getenv(name)), 64)
	if err != nil {
		return fallback
	}
	return value
}

// GrantAdoDiscordSignupCredit pays a new account its one signup credit. A problem
// here never stops the sign-up: the worst case is an account with no free credit.
// Without ADO_FINGERPRINT_KEY it pays nothing, because repeats couldn't be caught.
func GrantAdoDiscordSignupCredit(userId int, discordID string) {
	usd := adoEnvFloat("ADO_SIGNUP_GRANT_USD", adoSignupGrantDefaultUSD)
	limit := int(adoEnvFloat("ADO_SIGNUP_GRANT_DAILY_LIMIT", adoSignupGrantDefaultLimit))
	if usd <= 0 || limit <= 0 {
		return
	}
	key := strings.TrimSpace(os.Getenv("ADO_FINGERPRINT_KEY"))
	if key == "" {
		common.SysError("ADO_FINGERPRINT_KEY is not set, paying no signup credit")
		return
	}
	if !adoDiscordAccountOldEnough(discordID, time.Now()) {
		return
	}

	quota := common.QuotaRound(usd * common.QuotaPerUnit)
	fingerprint := common.GenerateHMACWithKey([]byte(key), discordID)
	outcome, err := model.PayAdoSignupGrant(userId, fingerprint, quota, limit)
	if err != nil {
		common.SysError("signup credit failed: " + err.Error())
		return
	}
	switch outcome {
	case model.AdoSignupGrantPaid:
		model.RecordLog(userId, model.LogTypeSystem, common.NewMessage("New user sign-up bonus {{quota}}", map[string]any{"quota": quota}))
	case model.AdoSignupGrantCapReached:
		adoAlertSignupGrantCap(limit)
	}
}

func adoAlertSignupGrantCap(limit int) {
	now := time.Now().Unix()
	last := adoLastCapAlertUnix.Load()
	if now-last < adoAlertCooldownSeconds || !adoLastCapAlertUnix.CompareAndSwap(last, now) {
		return
	}
	message := fmt.Sprintf("AdoRouter: signup credit paused, %d grants paid in the last 24 hours. New accounts get no credit until it drops. This may be abuse, or you may want to raise ADO_SIGNUP_GRANT_DAILY_LIMIT.", limit)
	common.SysError(message)
	webhook := strings.TrimSpace(os.Getenv("ADO_ALERT_WEBHOOK_URL"))
	if webhook == "" {
		return
	}
	go func() {
		body, err := common.Marshal(map[string]string{"content": message})
		if err != nil {
			return
		}
		client := http.Client{Timeout: 5 * time.Second}
		res, err := client.Post(webhook, "application/json", bytes.NewReader(body))
		if err != nil {
			common.SysError("alert webhook failed: " + err.Error())
			return
		}
		defer res.Body.Close()
		if res.StatusCode >= 300 {
			common.SysError(fmt.Sprintf("alert webhook failed: status %d", res.StatusCode))
		}
	}()
}
