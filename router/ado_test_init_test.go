package router

import "github.com/QuantumNous/new-api/middleware"

// Upstream's route tests exercise routes adorouter keeps closed in production.
func init() {
	middleware.AdoTextOnlyRelay = false
}
