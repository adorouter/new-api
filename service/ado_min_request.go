package service

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
)

// AdoRouter: the marketplace bills each request at least $0.00001, however few
// tokens it used, so a token price alone sells tiny requests on cheap lanes
// below cost. Its lanes are the groups new-api-sync names m1-<channel>-<model>.
const (
	adoMinRequestGroupPrefix = "m1-"
	adoMinRequestDefaultUSD  = 0.00001
)

// adoMinRequestQuota is the least a billable request routed through group may
// cost, or 0 when no minimum applies. ADO_MIN_REQUEST_USD overrides the amount;
// 0 turns it off.
func adoMinRequestQuota(group string) int {
	if !strings.HasPrefix(group, adoMinRequestGroupPrefix) {
		return 0
	}
	usd := adoEnvFloat("ADO_MIN_REQUEST_USD", adoMinRequestDefaultUSD)
	if usd <= 0 || usd > 1 {
		return 0
	}
	return common.QuotaRound(usd * common.QuotaPerUnit)
}
