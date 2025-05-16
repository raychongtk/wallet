package util

import "github.com/google/uuid"

// GetAssetAccount hardcode chart accounts for simplicity
func GetAssetAccount() (accountId uuid.UUID) {
	accountId, err := uuid.Parse("c19f00b4-c457-43f5-9e30-d10ada02a94f")
	if err != nil {
		return
	}
	return accountId
}

func GetLiabilityAccount() (accountId uuid.UUID) {
	accountId, err := uuid.Parse("d3b07384-d9a0-4f3b-8a2b-6c9e5b8b8f3c")
	if err != nil {
		return
	}
	return accountId
}
