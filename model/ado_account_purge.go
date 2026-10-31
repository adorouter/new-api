package model

import (
	"errors"
	"strconv"

	"github.com/QuantumNous/new-api/common"
	"gorm.io/gorm"
)

var (
	ErrAdoPurgeActiveSubscription = errors.New("cancel your subscription before deleting your account")
)

// PurgeUserAccount is what "delete my account" means on adorouter: the account
// and its personal data go for good, including the Discord ID and email, instead
// of upstream's soft delete. Kept on purpose: payment records (top_ups,
// subscription_orders) because tax law needs them, redeemed codes because they
// are other people's money, and audit rows with the identifying fields blanked.
func PurgeUserAccount(id int) error {
	if id == 0 {
		return errors.New("id is empty")
	}

	var role int
	if err := DB.Model(&User{}).Where("id = ?", id).Select("role").Scan(&role).Error; err != nil {
		return err
	}
	if role == common.RoleRootUser {
		return ErrCannotDeleteRootUser
	}
	active, err := HasActiveUserSubscription(id)
	if err != nil {
		return err
	}
	if active {
		return ErrAdoPurgeActiveSubscription
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		for _, personalData := range []any{
			&Task{},
			&Midjourney{},
			&Checkin{},
			&FreeConsent{},
			&ProviderConsent{},
			&TermsAcceptance{},
			&QuotaData{},
		} {
			if err := tx.Unscoped().Where("user_id = ?", id).Delete(personalData).Error; err != nil {
				return err
			}
		}
		return tx.Unscoped().Where("user_id = ?", strconv.Itoa(id)).Delete(&OAuthGrant{}).Error
	})
	if err != nil {
		return err
	}

	if _, err := HardDeleteUserById(id); err != nil {
		return err
	}

	// The log database can be separate, so these can't join the transaction. They
	// run last: a failure here leaves metadata that retention cleans up, never an
	// account that looks deleted but isn't.
	if err := LOG_DB.Where("user_id = ?", id).Delete(&Log{}).Error; err != nil {
		return err
	}
	return LOG_DB.Model(&AuditLog{}).Where("user_id = ?", id).
		Updates(map[string]any{"username": "", "ip": "", "user_agent": ""}).Error
}
