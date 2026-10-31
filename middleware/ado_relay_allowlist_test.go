package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdoTextOnlyRelay(t *testing.T) {
	gin.SetMode(gin.TestMode)
	old := AdoTextOnlyRelay
	AdoTextOnlyRelay = true
	t.Cleanup(func() { AdoTextOnlyRelay = old })

	tests := []struct {
		path string
		want int
	}{
		{"/v1/chat/completions", http.StatusOK},
		{"/v1/completions", http.StatusOK},
		{"/v1/responses", http.StatusOK},
		{"/v1/responses/compact", http.StatusOK},
		{"/v1/messages", http.StatusOK},
		{"/v1/embeddings", http.StatusOK},
		{"/v1/models", http.StatusOK},
		{"/v1/models/deepseek-flash", http.StatusOK},
		{"/v1beta/models/gemini:generateContent", http.StatusOK},
		{"/pg/chat/completions", http.StatusOK},
		{"/v1/images/generations", http.StatusNotFound},
		{"/v1/images/edits", http.StatusNotFound},
		{"/v1/videos", http.StatusNotFound},
		{"/v1/video/generations", http.StatusNotFound},
		{"/v1/tasks/some-plugin", http.StatusNotFound},
		{"/v1/audio/speech", http.StatusNotFound},
		{"/v1/realtime", http.StatusNotFound},
		{"/v1/alpha/search", http.StatusNotFound},
		{"/v1/media/some-token", http.StatusNotFound},
		{"/mj/submit/imagine", http.StatusNotFound},
		{"/fast/mj/submit/imagine", http.StatusNotFound},
		{"/suno/submit/music", http.StatusNotFound},
		{"/v1/chat/completions/extra", http.StatusNotFound},
		{"/v1/brand-new-upstream-route", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			engine := gin.New()
			engine.Use(RouteTag("relay"))
			engine.NoRoute(func(c *gin.Context) { c.Status(http.StatusOK) })

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			engine.ServeHTTP(rec, req)

			assert.Equal(t, tt.want, rec.Code)
			if tt.want == http.StatusNotFound {
				assert.Contains(t, rec.Body.String(), "ENDPOINT_NOT_AVAILABLE")
			}
		})
	}
}

func TestAdoTextOnlyRelayLeavesOtherRoutesAlone(t *testing.T) {
	gin.SetMode(gin.TestMode)
	old := AdoTextOnlyRelay
	AdoTextOnlyRelay = true
	t.Cleanup(func() { AdoTextOnlyRelay = old })

	engine := gin.New()
	engine.Use(RouteTag("api"))
	engine.NoRoute(func(c *gin.Context) { c.Status(http.StatusOK) })

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/status", nil))

	assert.Equal(t, http.StatusOK, rec.Code)
}
