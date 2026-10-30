package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createAdoPurgeUser(t *testing.T, username string, role int) User {
	t.Helper()
	user := User{
		Username:    username,
		Password:    "unused-password-hash",
		Role:        role,
		Status:      common.UserStatusEnabled,
		Group:       "default",
		AuthVersion: 1,
		AffCode:     username + "-aff",
		Email:       username + "@example.com",
		DiscordId:   username + "-discord",
		RegisterIp:  "203.0.113.7",
	}
	require.NoError(t, DB.Create(&user).Error)
	return user
}

func setupAdoPurgeTest(t *testing.T) {
	t.Helper()
	setupUserUpdateTestState(t)
	require.NoError(t, DB.AutoMigrate(&Midjourney{}, &Checkin{}, &FreeConsent{}, &ProviderConsent{}, &TermsAcceptance{}, &QuotaData{}, &OAuthGrant{}, &TopUp{}, &UserSubscription{}))
	require.NoError(t, LOG_DB.AutoMigrate(&AuditLog{}, &Log{}))
	t.Cleanup(func() {
		DB.Exec("DELETE FROM top_ups")
		DB.Exec("DELETE FROM user_subscriptions")
		LOG_DB.Exec("DELETE FROM audit_logs")
	})
}

func TestPurgeUserAccountRemovesPersonalDataAndLeavesOtherUsersAlone(t *testing.T) {
	setupAdoPurgeTest(t)
	gone := createAdoPurgeUser(t, "purge-gone", common.RoleCommonUser)
	kept := createAdoPurgeUser(t, "purge-kept", common.RoleCommonUser)

	for _, id := range []int{gone.Id, kept.Id} {
		require.NoError(t, DB.Create(&Token{UserId: id, Name: "key", Key: "key-for-user-" + common.Interface2String(id)}).Error)
		require.NoError(t, LOG_DB.Create(&Log{UserId: id, Username: "someone", Type: LogTypeConsume, Ip: "203.0.113.7"}).Error)
		require.NoError(t, DB.Create(&TopUp{UserId: id, Amount: 5, TradeNo: "trade-" + common.Interface2String(id), Status: "success"}).Error)
		require.NoError(t, LOG_DB.Create(&AuditLog{UserId: id, EventId: "event-" + common.Interface2String(id), Username: "someone", Ip: "203.0.113.7", UserAgent: "agent"}).Error)
	}

	require.NoError(t, PurgeUserAccount(gone.Id))

	var count int64
	require.NoError(t, DB.Unscoped().Model(&User{}).Where("id = ?", gone.Id).Count(&count).Error)
	assert.Zero(t, count, "the row must be gone, not soft deleted")
	require.NoError(t, DB.Unscoped().Model(&User{}).Where("discord_id = ?", gone.DiscordId).Count(&count).Error)
	assert.Zero(t, count, "the Discord ID must not survive")
	require.NoError(t, DB.Unscoped().Model(&User{}).Where("email = ?", gone.Email).Count(&count).Error)
	assert.Zero(t, count, "the email must not survive")
	require.NoError(t, DB.Unscoped().Model(&Token{}).Where("user_id = ?", gone.Id).Count(&count).Error)
	assert.Zero(t, count, "API keys must be gone")
	require.NoError(t, LOG_DB.Model(&Log{}).Where("user_id = ?", gone.Id).Count(&count).Error)
	assert.Zero(t, count, "usage logs must be gone")

	require.NoError(t, DB.Model(&TopUp{}).Where("user_id = ?", gone.Id).Count(&count).Error)
	assert.Equal(t, int64(1), count, "payment records stay for tax")

	var audit AuditLog
	require.NoError(t, LOG_DB.Where("user_id = ?", gone.Id).First(&audit).Error)
	assert.Empty(t, audit.Username)
	assert.Empty(t, audit.Ip)
	assert.Empty(t, audit.UserAgent)

	var other User
	require.NoError(t, DB.First(&other, kept.Id).Error)
	assert.Equal(t, kept.DiscordId, other.DiscordId)
	require.NoError(t, DB.Model(&Token{}).Where("user_id = ?", kept.Id).Count(&count).Error)
	assert.Equal(t, int64(1), count)
	require.NoError(t, LOG_DB.Model(&Log{}).Where("user_id = ?", kept.Id).Count(&count).Error)
	assert.Equal(t, int64(1), count)
	var otherAudit AuditLog
	require.NoError(t, LOG_DB.Where("user_id = ?", kept.Id).First(&otherAudit).Error)
	assert.Equal(t, "someone", otherAudit.Username)
}

func TestPurgeUserAccountRefusals(t *testing.T) {
	setupAdoPurgeTest(t)

	root := createAdoPurgeUser(t, "purge-root", common.RoleRootUser)
	require.ErrorIs(t, PurgeUserAccount(root.Id), ErrCannotDeleteRootUser)

	subscribed := createAdoPurgeUser(t, "purge-subscribed", common.RoleCommonUser)
	require.NoError(t, DB.Create(&UserSubscription{
		UserId:    subscribed.Id,
		PlanId:    1,
		Status:    "active",
		StartTime: common.GetTimestamp() - 60,
		EndTime:   common.GetTimestamp() + 3600,
	}).Error)
	require.ErrorIs(t, PurgeUserAccount(subscribed.Id), ErrAdoPurgeActiveSubscription)

	var count int64
	require.NoError(t, DB.Model(&User{}).Where("id IN ?", []int{root.Id, subscribed.Id}).Count(&count).Error)
	assert.Equal(t, int64(2), count)
}
