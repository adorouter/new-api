package model

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The public status page counts the models the catalog shows: a "[1m]"
// long-context alias is the same model, not one of its own.
func TestPublicStatusComponentsLeaveOutLongContextAliases(t *testing.T) {
	require.NoError(t, DB.AutoMigrate(&ModelStatusComponent{}))
	require.NoError(t, DB.Exec("DELETE FROM model_status_components").Error)
	t.Cleanup(func() { DB.Exec("DELETE FROM model_status_components") })

	now := common.GetTimestamp()
	for _, name := range []string{"claude-opus-5", "claude-opus-5[1m]", "deepseek-v4-pro"} {
		require.NoError(t, DB.Create(&ModelStatusComponent{ModelName: name, Public: true, CreatedAt: now, LastUpAt: now}).Error)
	}

	rows, err := GetAllPublicModelStatusComponents()
	require.NoError(t, err)
	names := make([]string, 0, len(rows))
	for _, row := range rows {
		names = append(names, row.ModelName)
	}
	assert.ElementsMatch(t, []string{"claude-opus-5", "deepseek-v4-pro"}, names)
}
