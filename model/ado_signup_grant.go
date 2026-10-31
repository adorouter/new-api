package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AdoSignupGrant remembers that a Discord account already had its one signup
// credit. It holds only a keyed one-way fingerprint and stays after the account is
// deleted, otherwise deleting and re-creating an account would pay out again.
type AdoSignupGrant struct {
	Fingerprint string `gorm:"primaryKey;type:varchar(64)"`
	CreatedAt   int64  `gorm:"bigint;index"`
}

func (AdoSignupGrant) TableName() string {
	return "ado_signup_grants"
}

type AdoSignupGrantOutcome string

const (
	AdoSignupGrantPaid        AdoSignupGrantOutcome = "paid"
	AdoSignupGrantAlreadyPaid AdoSignupGrantOutcome = "already_paid"
	AdoSignupGrantCapReached  AdoSignupGrantOutcome = "cap_reached"
)

// PayAdoSignupGrant credits quota once per fingerprint, and at most dailyLimit
// times in any 24 hours. Two sign-ups at the very same moment can pass the cap
// together, which costs one extra grant at worst, so it isn't worth a table lock.
func PayAdoSignupGrant(userId int, fingerprint string, quota int, dailyLimit int) (AdoSignupGrantOutcome, error) {
	outcome := AdoSignupGrantPaid
	err := DB.Transaction(func(tx *gorm.DB) error {
		var recent int64
		since := time.Now().Add(-24 * time.Hour).Unix()
		if err := tx.Model(&AdoSignupGrant{}).Where("created_at >= ?", since).Count(&recent).Error; err != nil {
			return err
		}
		if recent >= int64(dailyLimit) {
			outcome = AdoSignupGrantCapReached
			return nil
		}

		insert := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&AdoSignupGrant{Fingerprint: fingerprint, CreatedAt: time.Now().Unix()})
		if insert.Error != nil {
			return insert.Error
		}
		if insert.RowsAffected == 0 {
			outcome = AdoSignupGrantAlreadyPaid
			return nil
		}
		return tx.Model(&User{}).Where("id = ?", userId).Update("quota", gorm.Expr("quota + ?", quota)).Error
	})
	if err != nil {
		return "", err
	}
	if outcome == AdoSignupGrantPaid {
		if err := invalidateUserCache(userId); err != nil {
			return outcome, err
		}
	}
	return outcome, nil
}
