package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidSigner                     = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrTransferTxNotFound                = errors.Register(ModuleName, 1, "could not find transfer tx")
	ErrMintOrBurnNotFound                = errors.Register(ModuleName, 2, "could not find mint or burn message")
	ErrInvalidCommitment                 = errors.Register(ModuleName, 3, "invalid commitment")
	ErrClearingAccountNotFound           = errors.Register(ModuleName, 4, "could not find clearing account")
	ErrNotEnoughBalanceInPublic          = errors.Register(ModuleName, 5, "not enough balance in public pool")
	ErrReserveNotFound                   = errors.Register(ModuleName, 6, "could not find reserve")
	ErrNotEnoughBalanceInPrivate         = errors.Register(ModuleName, 7, "not enough balance in private pool")
	ErrNotEnoughUserBalanceInReserves    = errors.Register(ModuleName, 8, "not enough user balance in reserves")
	ErrNotEnoughUserBalanceInBank        = errors.Register(ModuleName, 9, "not enough user balance in bank")
	ErrInvalidTwilightAddress            = errors.Register(ModuleName, 10, "invalid twilight address")
	ErrInvalidInput                      = errors.Register(ModuleName, 11, "invalid input")
	ErrBankBalanceNotEqualReserveBalance = errors.Register(ModuleName, 12, "bank balance not equal to reserve balance")
	ErrDuplicateQqAccount                = errors.Register(ModuleName, 13, "duplicate qq account")
)
