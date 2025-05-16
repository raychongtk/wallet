package util

import "github.com/google/uuid"

// GetAssetAccount hardcode chart accounts for simplicity
func GetAssetAccount() (accountId uuid.UUID) {
	accountId, err := uuid.Parse("338b3f97-e428-4bff-9775-f759b5fccc4d")
	if err != nil {
		return
	}
	return accountId
}

func GetLiabilityAccount() (accountId uuid.UUID) {
	accountId, err := uuid.Parse("141e3fd8-c350-4b44-a2d5-2e2602aca72a")
	if err != nil {
		return
	}
	return accountId
}
