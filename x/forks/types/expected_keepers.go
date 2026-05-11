package types

import (
	"context"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetSupply(ctx context.Context, denom string) sdk.Coin
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
}

// StakingKeeper defines the expected staking keeper methods
type StakingKeeper interface {
	GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error)
	GetLastValidatorPower(ctx context.Context, operator sdk.ValAddress) (int64, error)
	GetLastTotalPower(ctx context.Context) (math.Int, error)
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	GetAllValidators(ctx context.Context) ([]stakingtypes.Validator, error)
	Validator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.ValidatorI, error)
}

// VoltKeeper defines the expected interface for the volt module
type VoltKeeper interface {
	UpdateBtcReserveAfterMint(ctx context.Context, mintedValue uint64, twilightAddress sdk.AccAddress, reserveAddress string) error
	UpdateBtcReserveAfterSweepProposal(ctx context.Context, reserveId uint64, reserveAddress string, judgeAddress string, btcBlockNumber uint64, btcRelayCapacityValue uint64, btcTxHash string, unlockHeight uint64, roundId uint64) error
	ConfirmWithdrawRequestsAfterSweepConfirmation(ctx context.Context, reserveId uint64, roundId uint64) error
	PruneReserveWithdrawSnapshot(ctx context.Context, reserveId uint64, roundId uint64)
	PruneRefundTxSnapshot(ctx context.Context, reserveId uint64, roundId uint64)
}
