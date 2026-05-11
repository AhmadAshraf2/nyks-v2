package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"twilight-project/nyks/x/bridge/types"
)

func (k Keeper) SetReserveAddressForJudge(ctx context.Context, judgeAddress sdk.AccAddress, reserveScript types.BtcScript, reserveAddress types.BtcAddress) {
	store := k.KVStore(ctx)
	regRes := &types.MsgRegisterReserveAddress{
		ReserveScript:  reserveScript.BtcScript,
		ReserveAddress: reserveAddress.BtcAddress,
		JudgeAddress:   judgeAddress.String(),
	}
	aKey := types.GetBtcRegisterReserveAddressKey(judgeAddress, reserveAddress)
	store.Set(aKey, k.cdc.MustMarshal(regRes))
}

func (k Keeper) IterateBtcReserveAddresses(ctx context.Context, cb func([]byte, types.MsgRegisterReserveAddress) bool) {
	store := k.KVStore(ctx)
	prefix := types.BtcReserveAddressKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgRegisterReserveAddress
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if cb(iter.Key(), res) {
			return
		}
	}
}

func (k Keeper) SetJudgeAddressForValidatorAddress(ctx context.Context, judgeAddress sdk.AccAddress, numOfSigners uint64, threshold uint64, signerApplicationFee uint64, arbitraryData string, validatorAddress sdk.ValAddress) error {
	store := k.KVStore(ctx)
	regJudge := &types.MsgBootstrapFragment{
		JudgeAddress:         judgeAddress.String(),
		NumOfSigners:         numOfSigners,
		Threshold:            threshold,
		SignerApplicationFee: signerApplicationFee,
		ArbitraryData:        arbitraryData,
		ValidatorAddress:     validatorAddress.String(),
	}
	aKey := types.GetBootstrapFragmentAddressKey(validatorAddress)
	store.Set(aKey, k.cdc.MustMarshal(regJudge))
	return nil
}

func (k Keeper) GetJudgeAddressForValidatorAddress(ctx context.Context, validatorAddress sdk.ValAddress) (sdk.AccAddress, error) {
	store := k.KVStore(ctx)
	prefix := types.JudgeAddressKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgBootstrapFragment
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if res.ValidatorAddress == validatorAddress.String() {
			judgeAddress, err := sdk.AccAddressFromBech32(res.JudgeAddress)
			if err != nil {
				return nil, err
			}
			return judgeAddress, nil
		}
	}
	return nil, types.ErrValidatorAddressNotFound.Wrapf("validator address %v", validatorAddress)
}

func (k Keeper) GetValidatorAddressForJudgeAddress(ctx context.Context, judgeAddress sdk.AccAddress) (sdk.ValAddress, error) {
	store := k.KVStore(ctx)
	prefix := types.JudgeAddressKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgBootstrapFragment
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if res.JudgeAddress == judgeAddress.String() {
			validatorAddress, err := sdk.ValAddressFromBech32(res.ValidatorAddress)
			if err != nil {
				return nil, err
			}
			return validatorAddress, nil
		}
	}
	return nil, types.ErrValidatorAddressNotFound.Wrapf("judge address %v", judgeAddress)
}

func (k Keeper) CheckJudgeValidatorInSet(ctx context.Context, judgeAddress sdk.AccAddress) bool {
	validatorAddress, err := k.GetValidatorAddressForJudgeAddress(ctx, judgeAddress)
	if err != nil {
		return false
	}
	_, err = k.StakingKeeper.GetValidator(ctx, validatorAddress)
	return err == nil
}

func (k Keeper) IterateRegisteredJudges(ctx context.Context, cb func([]byte, types.MsgBootstrapFragment) bool) {
	store := k.KVStore(ctx)
	prefix := types.JudgeAddressKey
	iter := store.Iterator(prefixRange(prefix))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var res types.MsgBootstrapFragment
		k.cdc.MustUnmarshal(iter.Value(), &res)
		if cb(iter.Key(), res) {
			return
		}
	}
}

func (k Keeper) SetBtcSignRefundMsg(ctx context.Context, signerAddress sdk.AccAddress, reserveId uint64, roundId uint64, singerPublicKey string, refundSignatures []string) error {
	store := k.KVStore(ctx)
	aKey := types.GetBtcSignRefundMsgKey(reserveId, roundId, signerAddress)
	signRefund := &types.MsgSignRefund{
		ReserveId:       reserveId,
		RoundId:         roundId,
		SignerPublicKey: singerPublicKey,
		RefundSignature: refundSignatures,
		SignerAddress:   signerAddress.String(),
	}
	store.Set(aKey, k.cdc.MustMarshal(signRefund))
	return nil
}

func (k Keeper) GetBtcSignRefundMsg(ctx context.Context, reserveId uint64, roundId uint64) ([]types.MsgSignRefund, bool) {
	store := k.KVStore(ctx)
	prefix := types.GetBtcSignRefundMsgPrefix(reserveId, roundId)
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	var msgs []types.MsgSignRefund
	for ; iterator.Valid(); iterator.Next() {
		var msg types.MsgSignRefund
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}
	if len(msgs) == 0 {
		return nil, false
	}
	return msgs, true
}

func (k Keeper) GetBtcSignRefundMsgWithOracleAddress(ctx context.Context, reserveId uint64, roundId uint64, signerAddress sdk.AccAddress) (*types.MsgSignRefund, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetBtcSignRefundMsgKey(reserveId, roundId, signerAddress)
	if !store.Has(aKey) {
		return nil, false
	}
	bz := store.Get(aKey)
	var signRefund types.MsgSignRefund
	k.cdc.MustUnmarshal(bz, &signRefund)
	return &signRefund, true
}

func (k Keeper) GetBtcSignSweepMsgWithOracleAddress(ctx context.Context, reserveId uint64, roundId uint64, signerAddress sdk.AccAddress) (*types.MsgSignSweep, bool) {
	store := k.KVStore(ctx)
	key := types.GetBtcSignSweepMsgKey(reserveId, roundId, signerAddress)
	if !store.Has(key) {
		return nil, false
	}
	bz := store.Get(key)
	var signSweep types.MsgSignSweep
	k.cdc.MustUnmarshal(bz, &signSweep)
	return &signSweep, true
}

func (k Keeper) SetBtcSignSweepMsg(ctx context.Context, signerAddress sdk.AccAddress, reserveId uint64, roundId uint64, singerPublicKey string, sweepSignatures []string) error {
	store := k.KVStore(ctx)
	aKey := types.GetBtcSignSweepMsgKey(reserveId, roundId, signerAddress)
	signSweep := &types.MsgSignSweep{
		ReserveId:       reserveId,
		RoundId:         roundId,
		SignerPublicKey: singerPublicKey,
		SweepSignature:  sweepSignatures,
		SignerAddress:   signerAddress.String(),
	}
	store.Set(aKey, k.cdc.MustMarshal(signSweep))
	return nil
}

func (k Keeper) GetBtcSignSweepMsg(ctx context.Context, reserveId uint64, roundId uint64) ([]types.MsgSignSweep, bool) {
	store := k.KVStore(ctx)
	prefix := types.GetBtcSignSweepMsgPrefix(reserveId, roundId)
	iterator := storetypes.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()
	var msgs []types.MsgSignSweep
	for ; iterator.Valid(); iterator.Next() {
		var msg types.MsgSignSweep
		k.cdc.MustUnmarshal(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}
	if len(msgs) == 0 {
		return nil, false
	}
	return msgs, true
}

func (k Keeper) SetBtcBroadcastTxSweepMsg(ctx context.Context, reserveId uint64, roundId uint64, judgeAddress sdk.AccAddress, signedSweepTx string) error {
	store := k.KVStore(ctx)
	aKey := types.GetBtcBroadcastTxSweepMsgKey(reserveId, roundId)
	msg := &types.MsgBroadcastTxSweep{
		ReserveId:     reserveId,
		RoundId:       roundId,
		SignedSweepTx: signedSweepTx,
		JudgeAddress:  judgeAddress.String(),
	}
	store.Set(aKey, k.cdc.MustMarshal(msg))
	return nil
}

func (k Keeper) GetBtcBroadcastTxSweepMsg(ctx context.Context, reserveId uint64, roundId uint64) (*types.MsgBroadcastTxSweep, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetBtcBroadcastTxSweepMsgKey(reserveId, roundId)
	if !store.Has(aKey) {
		return nil, false
	}
	bz := store.Get(aKey)
	var msg types.MsgBroadcastTxSweep
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}

func (k Keeper) SetBtcBroadcastTxRefundMsg(ctx context.Context, reserveId uint64, roundId uint64, judgeAddress sdk.AccAddress, signedRefundTx string) error {
	store := k.KVStore(ctx)
	aKey := types.GetBtcBroadcastTxRefundMsgKey(reserveId, roundId)
	msg := &types.MsgBroadcastTxRefund{
		ReserveId:      reserveId,
		RoundId:        roundId,
		SignedRefundTx: signedRefundTx,
		JudgeAddress:   judgeAddress.String(),
	}
	store.Set(aKey, k.cdc.MustMarshal(msg))
	return nil
}

func (k Keeper) GetBtcBroadcastTxRefundMsg(ctx context.Context, reserveId uint64, roundId uint64) (*types.MsgBroadcastTxRefund, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetBtcBroadcastTxRefundMsgKey(reserveId, roundId)
	if !store.Has(aKey) {
		return nil, false
	}
	bz := store.Get(aKey)
	var msg types.MsgBroadcastTxRefund
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}

func (k Keeper) SetBtcProposeRefundHashMsg(ctx context.Context, judgeAddress sdk.AccAddress, refundHash string) error {
	store := k.KVStore(ctx)
	aKey := types.GetBtcProposeRefundHashMsgKey(judgeAddress, refundHash)
	msg := &types.MsgProposeRefundHash{
		JudgeAddress: judgeAddress.String(),
		RefundHash:   refundHash,
	}
	store.Set(aKey, k.cdc.MustMarshal(msg))
	return nil
}

func (k Keeper) SetUnsignedTxSweepMsg(ctx context.Context, txId string, unsignedSweepTx string, reserveId uint64, roundId uint64, judgeAddress sdk.AccAddress) error {
	store := k.KVStore(ctx)
	aKey := types.GetUnsignedTxSweepMsgKey(reserveId, roundId)
	msg := &types.MsgUnsignedTxSweep{
		TxId:               txId,
		BtcUnsignedSweepTx: unsignedSweepTx,
		ReserveId:          reserveId,
		RoundId:            roundId,
		JudgeAddress:       judgeAddress.String(),
	}
	store.Set(aKey, k.cdc.MustMarshal(msg))
	return nil
}

func (k Keeper) GetUnsignedTxSweepMsg(ctx context.Context, reserveId uint64, roundId uint64) (*types.MsgUnsignedTxSweep, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetUnsignedTxSweepMsgKey(reserveId, roundId)
	if !store.Has(aKey) {
		return nil, false
	}
	bz := store.Get(aKey)
	var msg types.MsgUnsignedTxSweep
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}

func (k Keeper) SetUnsignedTxRefundMsg(ctx context.Context, reserveId uint64, roundId uint64, btcUnsignedRefundTx string, judgeAddress sdk.AccAddress) error {
	store := k.KVStore(ctx)
	aKey := types.GetUnsignedTxRefundMsgKey(reserveId, roundId)
	msg := &types.MsgUnsignedTxRefund{
		ReserveId:           reserveId,
		RoundId:             roundId,
		BtcUnsignedRefundTx: btcUnsignedRefundTx,
		JudgeAddress:        judgeAddress.String(),
	}
	store.Set(aKey, k.cdc.MustMarshal(msg))
	return nil
}

func (k Keeper) GetUnsignedTxRefundMsg(ctx context.Context, reserveId uint64, roundId uint64) (*types.MsgUnsignedTxRefund, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetUnsignedTxRefundMsgKey(reserveId, roundId)
	if !store.Has(aKey) {
		return nil, false
	}
	bz := store.Get(aKey)
	var msg types.MsgUnsignedTxRefund
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}

func (k Keeper) SetProposeSweepAddress(ctx context.Context, btcAddress types.BtcAddress, btcScript string, reserveId uint64, roundId uint64, judgeAddress sdk.AccAddress) error {
	store := k.KVStore(ctx)
	aKey := types.GetProposeSweepAddressMsgKey(reserveId, roundId)
	msg := &types.MsgProposeSweepAddress{
		BtcAddress:   btcAddress.BtcAddress,
		BtcScript:    btcScript,
		ReserveId:    reserveId,
		RoundId:      roundId,
		JudgeAddress: judgeAddress.String(),
	}
	store.Set(aKey, k.cdc.MustMarshal(msg))
	return nil
}

func (k Keeper) GetProposeSweepAddress(ctx context.Context, reserveId uint64, roundId uint64) (*types.MsgProposeSweepAddress, bool) {
	store := k.KVStore(ctx)
	aKey := types.GetProposeSweepAddressMsgKey(reserveId, roundId)
	if !store.Has(aKey) {
		return nil, false
	}
	bz := store.Get(aKey)
	var msg types.MsgProposeSweepAddress
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}

func (k Keeper) LockProposeSweepAddress(ctx context.Context) {
	store := k.KVStore(ctx)
	store.Set(types.ProposeSweepAddressLockKey, []byte{1})
}

func (k Keeper) UnlockProposeSweepAddress(ctx context.Context) {
	store := k.KVStore(ctx)
	store.Delete(types.ProposeSweepAddressLockKey)
}

func (k Keeper) IsProposeSweepAddressLocked(ctx context.Context) bool {
	store := k.KVStore(ctx)
	return store.Has(types.ProposeSweepAddressLockKey)
}

func (k Keeper) GetAllProposedSweepAddresses(ctx context.Context, limit uint64) ([]types.MsgProposeSweepAddress, error) {
	store := k.KVStore(ctx)
	prefix := types.ProposeSweepAddressMsg
	iter := store.ReverseIterator(prefixRange(prefix))
	defer iter.Close()
	var results []types.MsgProposeSweepAddress
	var count uint64
	for ; iter.Valid() && count < limit; iter.Next() {
		var res types.MsgProposeSweepAddress
		k.cdc.MustUnmarshal(iter.Value(), &res)
		results = append(results, res)
		count++
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].ReserveId != results[j].ReserveId {
			return results[i].ReserveId > results[j].ReserveId
		}
		return results[i].RoundId > results[j].RoundId
	})
	return results, nil
}

func hashSignatures(signatures []string) string {
	concatenated := strings.Join(signatures, "|")
	hash := sha256.Sum256([]byte(concatenated))
	return hex.EncodeToString(hash[:])
}
