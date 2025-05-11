package service_service

import (
	"fmt"
	"nearbyassist/internal/utils"
)

func computeSignature(vendorId, title, description, pricingType string, hashFunc func([]byte) (string, error)) string {
	rawStr := fmt.Sprintf("%s_%s_%s_%s", vendorId, title, description, pricingType)
	signature := utils.Must(hashFunc([]byte(rawStr)))
	return signature
}
