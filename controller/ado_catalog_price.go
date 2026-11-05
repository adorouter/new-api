package controller

import (
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/pkg/billingexpr"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
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

// hasConfiguredPrice is false for a model an admin has not priced yet. The
// gateway lists such a model at its placeholder ratio (37.5, $75 per 1M input
// before the group ratio) and refuses every call to it with model_price_error,
// so the public catalog hides it until it has a price. Self-use mode bills the
// placeholder, so there the model really is callable at that price and stays.
func hasConfiguredPrice(m model.Pricing) bool {
	if m.QuotaType != 0 || m.BillingMode != "" {
		return true
	}
	_, ok, _ := ratio_setting.GetModelRatio(m.ModelName)
	return ok
}

func withConfiguredPrice(pricing []model.Pricing) []model.Pricing {
	out := make([]model.Pricing, 0, len(pricing))
	for _, m := range pricing {
		if hasConfiguredPrice(m) {
			out = append(out, m)
		}
	}
	return out
}
