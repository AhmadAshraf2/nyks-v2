package types

import (
	"context"

	"cosmossdk.io/core/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	forkstypes "twilight-project/nyks/x/forks/types"
	volttypes "twilight-project/nyks/x/volt/types"
)

type AccountKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
}

type NyksKeeper interface {
	CheckOrchestratorValidatorInSet(ctx context.Context, orchestrator string) (sdk.ValAddress, error)
	ClaimHandlerCommon(ctx sdk.Context, msgAny *codectypes.Any, valAddr sdk.ValAddress, msg forkstypes.BtcProposal) error
}

type StakingKeeper interface {
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
}

type VoltKeeper interface {
	GetBtcReserve(ctx context.Context, reserveId uint64) (*volttypes.BtcReserve, error)
	GetBtcReserveIdByAddress(ctx context.Context, reserveAddress string) (uint64, error)
	RegisterNewBtcReserve(ctx context.Context, judgeAddress sdk.AccAddress, reserveAddress string) (uint64, error)
	SetBtcDeposit(ctx context.Context, btcDepositAddress string, twilightAddress sdk.AccAddress, twilightStakingAmount uint64, btcSatoshiTestAmount uint64) error
	GetBtcDepositAddressByTwilightAddress(ctx context.Context, twilightAddress sdk.AccAddress) (btcDeposit *volttypes.BtcDepositAddress, found bool)
	GetClearingAccount(ctx context.Context, twilightAddress sdk.AccAddress) (*volttypes.ClearingAccount, bool)
	GetAllBtcRegisteredDepositAddresses(ctx context.Context) (btcDepositAddresses []volttypes.BtcDepositAddress)
	CheckBtcAddress(ctx context.Context, twilightAddress sdk.Address, btcAddress string, newSatoshiTestAmount uint64) bool
	SetBtcWithdrawRequest(ctx context.Context, twilightAddress sdk.AccAddress, reserveId uint64, withdrawAddress string, withdrawAmount uint64) (*uint32, error)
	CheckClearingAccountBalance(ctx context.Context, twilightAddress sdk.AccAddress, reserveId uint64, amount uint64) error
	CheckReserveWithdrawSnapshot(ctx context.Context, btcTxHex string, reserveId uint64, roundId uint64) (bool, error)
	CheckBtcReserveExists(ctx context.Context, reserveId uint64) bool
	RegisterNewFragment(ctx context.Context, judgeAddress sdk.AccAddress, threshold uint64, applicationFee uint64, numOfSigners uint64, fragmentFeeBips uint64, arbitraryData string) (uint64, error)
	SetNewSweepProposalReceived(ctx context.Context, reserveId uint64, roundId uint64)
}
