package aws

import (
	"errors"
	"time"
)

var ErrNotWhitelisted = errors.New("not whitelisted")

// InstanceIdentityDocument represents the AWS Instance Identity Document structure.
type InstanceIdentityDocument struct {
	DevpayProductCodes      []string  `json:"devpayProductCodes"`
	MarketplaceProductCodes []string  `json:"marketplaceProductCodes"`
	AvailabilityZone        string    `json:"availabilityZone"`
	PrivateIP               string    `json:"privateIp"`
	InstanceID              string    `json:"instanceId"`
	BillingProducts         []string  `json:"billingProducts"`
	AccountID               string    `json:"accountId"`
	ImageID                 string    `json:"imageId"`
	PendingTime             time.Time `json:"pendingTime"`
	Architecture            string    `json:"architecture"`
	KernelID                string    `json:"kernelId"`
	RamdiskID               string    `json:"ramdiskId"`
	Region                  string    `json:"region"`
}

func CheckWhitelist(value string, whitelist []string) bool {
	if len(whitelist) < 1 {
		return true
	}

	for _, i := range whitelist {
		if value == i {
			return true
		}
	}

	return false
}
