package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/forks module sentinel errors
var (
	ErrInvalidSigner           = errors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrInternal                = errors.Register(ModuleName, 1, "internal")
	ErrDuplicate               = errors.Register(ModuleName, 2, "duplicate")
	ErrInvalid                 = errors.Register(ModuleName, 3, "invalid")
	ErrTimeout                 = errors.Register(ModuleName, 4, "timeout")
	ErrResetDelegateKeys       = errors.Register(ModuleName, 5, "can not set btcOracle address mapping more than once")
	ErrInvalidBtcPublicKey     = errors.Register(ModuleName, 6, "invalid btc public key")
	ErrNonContiguousEventNonce = errors.Register(ModuleName, 9, "non contiguous event nonce")
	ErrInvalidValidator        = errors.Register(ModuleName, 10, "invalid validator")
	ErrAttestationOverflow = errors.Register(ModuleName, 11, "integer overflow in attestation")
)
