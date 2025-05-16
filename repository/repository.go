package repository

import "github.com/google/wire"

var (
	WireSet = wire.NewSet(
		ProvideUserRepository,
		ProvideMovementRepository,
		ProvideAccountRepository,
		ProvideWalletRepository,
		ProvideTransactionRepository,
	)
)
