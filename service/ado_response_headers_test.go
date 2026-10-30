package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShouldCopyUpstreamHeaderOnlyPassesTheAllowlist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)

	tests := []struct {
		header string
		want   bool
	}{
		{"Content-Type", true},
		{"content-encoding", true},
		{"Content-Disposition", true},
		{"Cache-Control", true},
		{"X-Codex-Turn-State", true},
		{"Content-Length", false},
		{"Set-Cookie", false},
		{"X-Ratelimit-Remaining-Tokens", false},
		{"Anthropic-Organization-Id", false},
		{"Openai-Organization", false},
		{"X-Some-New-Provider-Header", false},
		{common.RequestIdKey, false},
	}
	for _, tt := range tests {
		t.Run(tt.header, func(t *testing.T) {
			assert.Equal(t, tt.want, ShouldCopyUpstreamHeader(c, tt.header, []string{"value"}))
		})
	}
}
