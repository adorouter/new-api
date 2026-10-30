package service

import (
	"time"

	"github.com/QuantumNous/new-api/common"
)

const adoLogRetentionDays = 90

// StartAdoLogRetention keeps usage logs to 90 days by queueing upstream's log
// cleanup task once a day. Upstream only offers that cleanup as a manual action.
func StartAdoLogRetention() {
	if !common.IsMasterNode {
		return
	}
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for ; true; <-ticker.C {
			cutoff := time.Now().AddDate(0, 0, -adoLogRetentionDays).Unix()
			if _, err := StartLogCleanupTask(cutoff); err != nil {
				common.SysError("usage log retention failed to start: " + err.Error())
			}
		}
	}()
}
