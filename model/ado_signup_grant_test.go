package model

import (
	"testing"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdoSignupGrantTest(t *testing.T) {
	t.Helper()
	setupAdoPurgeTest(t)
	require.NoError(t, DB.AutoMigrate(&AdoSignupGrant{}))
	require.NoError(t, DB.Exec("DELETE FROM ado_signup_grants").Error)
	t.Cleanup(func() { DB.Exec("DELETE FROM ado_signup_grants") })
}

func adoUserQuota(t *testing.T, id int) int {
	t.Helper()
	var user User
	require.NoError(t, DB.First(&user, id).Error)
	return user.Quota
}

func TestPayAdoSignupGrantPaysOncePerDiscordAccountEvenAfterDeletion(t *testing.T) {
	setupAdoSignupGrantTest(t)
	first := createAdoPurgeUser(t, "grant-first", common.RoleCommonUser)

	outcome, err := PayAdoSignupGrant(first.Id, "fingerprint-a", 250000, 20)
	require.NoError(t, err)
	assert.Equal(t, AdoSignupGrantPaid, outcome)
	assert.Equal(t, 250000, adoUserQuota(t, first.Id))

	outcome, err = PayAdoSignupGrant(first.Id, "fingerprint-a", 250000, 20)
	require.NoError(t, err)
	assert.Equal(t, AdoSignupGrantAlreadyPaid, outcome)
	assert.Equal(t, 250000, adoUserQuota(t, first.Id), "a repeat must not add credit")

	require.NoError(t, PurgeUserAccount(first.Id))
	again := createAdoPurgeUser(t, "grant-again", common.RoleCommonUser)

	outcome, err = PayAdoSignupGrant(again.Id, "fingerprint-a", 250000, 20)
	require.NoError(t, err)
	assert.Equal(t, AdoSignupGrantAlreadyPaid, outcome, "deleting the account must not reset the grant")
	assert.Zero(t, adoUserQuota(t, again.Id))
}

func TestPayAdoSignupGrantStopsAtTheDailyCap(t *testing.T) {
	setupAdoSignupGrantTest(t)
	one := createAdoPurgeUser(t, "grant-one", common.RoleCommonUser)
	two := createAdoPurgeUser(t, "grant-two", common.RoleCommonUser)
	three := createAdoPurgeUser(t, "grant-three", common.RoleCommonUser)

	_, err := PayAdoSignupGrant(one.Id, "fingerprint-1", 250000, 2)
	require.NoError(t, err)
	_, err = PayAdoSignupGrant(two.Id, "fingerprint-2", 250000, 2)
	require.NoError(t, err)

	outcome, err := PayAdoSignupGrant(three.Id, "fingerprint-3", 250000, 2)
	require.NoError(t, err)
	assert.Equal(t, AdoSignupGrantCapReached, outcome)
	assert.Zero(t, adoUserQuota(t, three.Id))

	var count int64
	require.NoError(t, DB.Model(&AdoSignupGrant{}).Where("fingerprint = ?", "fingerprint-3").Count(&count).Error)
	assert.Zero(t, count, "a capped sign-up isn't marked as paid")

	require.NoError(t, DB.Model(&AdoSignupGrant{}).Where("1 = 1").Update("created_at", time.Now().Add(-25*time.Hour).Unix()).Error)
	outcome, err = PayAdoSignupGrant(three.Id, "fingerprint-3", 250000, 2)
	require.NoError(t, err)
	assert.Equal(t, AdoSignupGrantPaid, outcome, "the cap is a rolling 24 hours")
}
