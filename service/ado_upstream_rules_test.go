package service

import (
	"errors"
	"net/http"
	"testing"

	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/stretchr/testify/assert"
)

// Marketplace pinned-merchant errors fail over to a sibling lane and reach the
// customer in English, never as the marketplace's Chinese text with its request
// ids. A delisted model takes the lane out of rotation at once.
func TestMarketplaceMerchantErrorsFailOverInEnglish(t *testing.T) {
	cases := []struct {
		name       string
		message    string
		disableNow bool
		userText   string
	}{
		{
			name:       "model delisted by the merchant",
			message:    "固定商家当前不可用。 原因：该商家未上架或未开放本次请求的模型能力。 解决方案：请切换固定商家，或开启智能路由/兜底后重试。 request_id: 2026110813425637",
			disableNow: true,
			userText:   "This provider can no longer serve this model. We have removed it from rotation, so retry and another provider will take it.",
		},
		{
			name:     "any other pinned-merchant error",
			message:  "固定商家当前不可用。 原因：一个尚未见过的原因。 request_id: 2026110813425637",
			userText: "This provider could not take the request right now. Retry and another provider will pick it up.",
		},
		{
			name:     "a specific rule still wins over the catch-all",
			message:  "固定商家当前不可用。 原因：该商家上游正在限流或容量不足。 request_id: 2026110812303924",
			userText: "The providers for this model are rate limited right now. Nothing is used up on your side. Try again in a little while.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := types.NewOpenAIError(errors.New(tc.message), types.ErrorCodeBadResponse, http.StatusServiceUnavailable)
			class := ClassifyUpstreamError(err)
			assert.True(t, class.Known)
			assert.True(t, class.Failover)
			assert.Equal(t, tc.disableNow, class.DisableNow)
			text, ok := UpstreamUserMessage(err)
			assert.True(t, ok)
			assert.Equal(t, tc.userText, text)
			assert.NotContains(t, text, "request_id")
		})
	}
}
