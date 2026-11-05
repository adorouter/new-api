package controller

import (
	"encoding/json"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/gin-gonic/gin"
)

// adoKeepStoredMetadata fills in the stored metadata when an update request leaves
// the field out. Model.Update writes the metadata column unconditionally, and the
// console's metadata form never sent it, so saving a model's description erased
// its context window and max output tokens (which the chat clamps on). An
// explicit "metadata": "" still clears it.
func adoKeepStoredMetadata(c *gin.Context, m *model.Model) error {
	raw, _ := c.Get(gin.BodyBytesKey)
	body, _ := raw.([]byte)
	var fields map[string]json.RawMessage
	if err := common.Unmarshal(body, &fields); err != nil {
		return err
	}
	if _, sent := fields["metadata"]; sent {
		return nil
	}
	var stored model.Model
	if err := model.DB.Select("metadata").First(&stored, m.Id).Error; err != nil {
		return err
	}
	m.Metadata = stored.Metadata
	return nil
}
