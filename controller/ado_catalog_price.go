package controller

import (
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/pkg/billingexpr"
)

// adoExpressionPrices quotes a model billed by expression from that expression,
// in USD per million tokens. The ratio tables hold only a placeholder for these
// models, which listed DeepSeek at $75 per million. The quote is the price right
// now, so it follows peak and off-peak. Charges never read it.
func adoExpressionPrices(m model.Pricing) (input, output float64, ok bool) {
	if m.BillingMode != "tiered_expr" || m.BillingExpr == "" || billingexpr.UsesFixedPricing(m.BillingExpr) {
		return 0, 0, false
	}
	// Expression coefficients are USD per million tokens, so one token reads the rate.
	input, _, err := billingexpr.RunExpr(m.BillingExpr, billingexpr.TokenParams{P: 1})
	if err != nil {
		return 0, 0, false
	}
	output, _, err = billingexpr.RunExpr(m.BillingExpr, billingexpr.TokenParams{C: 1})
	if err != nil {
		return 0, 0, false
	}
	return input, output, true
}
