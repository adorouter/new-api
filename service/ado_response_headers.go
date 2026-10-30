package service

import "strings"

// adoForwardedResponseHeader is the allowlist of provider headers a client may
// receive. The rest describe the provider, not adorouter: cookies it would set on
// our domain, its rate limits, request IDs and account details.
func adoForwardedResponseHeader(name string) bool {
	switch strings.ToLower(name) {
	case "content-type", "content-encoding", "content-disposition", "content-language", "cache-control":
		return true
	// Protocol state for Codex-style clients on the Responses API. Dropping them
	// breaks those clients; they say nothing about the provider's account.
	case "x-reasoning-included", "x-codex-turn-state":
		return true
	}
	return false
}
