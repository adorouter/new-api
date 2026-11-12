package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The marketplace bills at least $0.00001 a request; its lanes must not settle
// under that, and no other group gets a minimum.
func TestAdoMinRequestQuota(t *testing.T) {
	for _, tc := range []struct {
		name  string
		group string
		env   string
		want  int
	}{
		{name: "marketplace lane, default $0.00001", group: "m1-6563-gpt-4o-mini", want: 5},
		{name: "marketplace lane, overridden", group: "m1-5052-glm-5.3-flash", env: "0.00002", want: 10},
		{name: "turned off", group: "m1-5052-glm-5.3-flash", env: "0", want: 0},
		{name: "absurd override ignored", group: "m1-5052-glm-5.3-flash", env: "50", want: 0},
		{name: "official lane", group: "default", want: 0},
		{name: "name merely containing m1", group: "xm1-1-model", want: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("ADO_MIN_REQUEST_USD", tc.env)
			assert.Equal(t, tc.want, adoMinRequestQuota(tc.group))
		})
	}
}
