package middleware

import (
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// adoConsolePublicPaths are upstream's marketing and password pages. adorouter's
// website replaces them, so a visitor to the API's address lands there instead.
// Sign-in, the OAuth callback and setup stay, because Discord sign-in passes through them.
var adoConsolePublicPaths = []string{
	"/",
	"/about",
	"/pricing",
	"/rankings",
	"/privacy-policy",
	"/user-agreement",
	"/sign-up",
	"/register",
	"/forgot-password",
	"/reset",
	"/user/reset",
}

// AdoRedirectPublicPages sends those pages to ADO_SITE_URL. Unset, as in local
// development, the console keeps all its pages.
func AdoRedirectPublicPages() gin.HandlerFunc {
	site := strings.TrimRight(strings.TrimSpace(os.Getenv("ADO_SITE_URL")), "/")
	return func(c *gin.Context) {
		path := strings.TrimRight(c.Request.URL.Path, "/")
		if path == "" {
			path = "/"
		}
		isPublicPage := slices.Contains(adoConsolePublicPaths, path) || strings.HasPrefix(path, "/pricing/")
		if site != "" && isPublicPage {
			c.Redirect(http.StatusFound, site)
			c.Abort()
			return
		}
		c.Next()
	}
}
