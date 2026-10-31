package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAdoRedirectPublicPages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name string
		site string
		path string
		want int
	}{
		{"home goes to the website", "https://adorouter.com/", "/", http.StatusFound},
		{"pricing goes to the website", "https://adorouter.com", "/pricing", http.StatusFound},
		{"model pricing goes to the website", "https://adorouter.com", "/pricing/deepseek-flash", http.StatusFound},
		{"trailing slash goes to the website", "https://adorouter.com", "/about/", http.StatusFound},
		{"discord callback stays", "https://adorouter.com", "/oauth/discord", http.StatusOK},
		{"admin sign-in stays", "https://adorouter.com", "/sign-in", http.StatusOK},
		{"first-run setup stays", "https://adorouter.com", "/setup", http.StatusOK},
		{"console pages stay", "https://adorouter.com", "/channels", http.StatusOK},
		{"unset keeps every page", "", "/pricing", http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ADO_SITE_URL", tt.site)
			engine := gin.New()
			engine.Use(AdoRedirectPublicPages())
			engine.NoRoute(func(c *gin.Context) { c.Status(http.StatusOK) })

			rec := httptest.NewRecorder()
			engine.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			assert.Equal(t, tt.want, rec.Code)
			if tt.want == http.StatusFound {
				assert.Equal(t, "https://adorouter.com", rec.Header().Get("Location"))
			}
		})
	}
}
