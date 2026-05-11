package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalidSigner                  = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInvalid                        = errors.Register(ModuleName, 1, "invalid")
	ErrDuplicate                      = errors.Register(ModuleName, 2, "duplicate")
	ErrResetBtcAddress                = errors.Register(ModuleName, 3, "can not set btc to twilight address mapping more than once")
	ErrInvalidBtcAddress              = errors.Register(ModuleName, 4, "invalid btc address")
	ErrInvalidTwilightAddress         = errors.Register(ModuleName, 5, "invalid twilight address")
	ErrJudgeAddressNotFound           = errors.Register(ModuleName, 6, "judge address not found")
	ErrValidatorAddressNotFound       = errors.Register(ModuleName, 7, "validator address not found")
	ErrJudgeValidatorNotFound         = errors.Register(ModuleName, 8, "validator for the judge not found")
	ErrInsufficientBalance            = errors.Register(ModuleName, 9, "insufficient user balance in reserve")
	ErrInsufficientBalanceInBank      = errors.Register(ModuleName, 10, "insufficient balance in bank")
	ErrClearingAccountDoesNotExist    = errors.Register(ModuleName, 11, "clearing account does not exist")
	ErrBtcAddressAlreadyExists        = errors.Register(ModuleName, 12, "btc address already exists")
	ErrProposeSweepAddressIsLocked    = errors.Register(ModuleName, 13, "propose sweep address is locked")
	ErrFragmentNotFound               = errors.Register(ModuleName, 14, "fragment not found")
	ErrMaxReservesPerFragmentExceeded = errors.Register(ModuleName, 15, "maximum reserves per fragment exceeded")
	ErrJudgeMismatch                  = errors.Register(ModuleName, 16, "judge mismatch")
)
