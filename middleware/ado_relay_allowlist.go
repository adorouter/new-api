package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdoTextOnlyRelay limits the relay to text endpoints. Image, video, music and
// task routes stay closed because the gateway stores their prompts, which would
// break adorouter's "we don't store your prompts" promise. It is a variable only
// so upstream's route tests can run with every route open.
var AdoTextOnlyRelay = true

var adoRelayExactPaths = []string{
	"/v1/chat/completions",
	"/v1/completions",
	"/v1/responses",
	"/v1/responses/compact",
	"/v1/messages",
	"/v1/messages/count_tokens",
	"/v1/embeddings",
	"/v1/rerank",
	"/v1/models",
	"/v1beta/models",
	"/v1beta/openai/models",
	"/pg/chat/completions",
}

var adoRelayPathPrefixes = []string{
	"/v1/models/",
	"/v1beta/models/",
	"/v1beta/openai/models/",
}

// adoRelayPathAllowed is an allowlist so a route added by an upstream merge
// starts closed instead of quietly storing prompts.
func adoRelayPathAllowed(path string) bool {
	if slices.Contains(adoRelayExactPaths, path) {
		return true
	}
	return slices.ContainsFunc(adoRelayPathPrefixes, func(prefix string) bool {
		return strings.HasPrefix(path, prefix)
	})
}

func abortAdoClosedRelayPath(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"error": gin.H{
			"message": "this endpoint isn't available: AdoRouter currently serves text models only",
			"type":    "invalid_request_error",
			"code":    "ENDPOINT_NOT_AVAILABLE",
		},
	})
}
