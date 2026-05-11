package keeper

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"twilight-project/nyks/x/forks/types"
)

func (k Keeper) Attest(
	ctx sdk.Context,
	proposal types.BtcProposal,
	valAddr sdk.ValAddress,
	anyProposal *codectypes.Any,
) (*types.Attestation, error) {
	hash, err := proposal.ProposalHash()
	if err != nil {
		return nil, fmt.Errorf("unable to compute claim hash: %w", err)
	}
	att := k.GetAttestation(ctx, proposal.GetHeight(), hash)

	if att == nil {
		att = &types.Attestation{
			Observed: false,
			Votes:    []string{},
			Height:   uint64(ctx.BlockHeight()),
			Proposal: anyProposal,
		}
	} else {
		for _, s := range att.Votes {
			if valAddr.String() == s {
				return nil, errors.New("duplicate vote")
			}
		}
	}

	att.Votes = append(att.Votes, valAddr.String())
	k.SetAttestation(ctx, proposal.GetHeight(), hash, att)
	k.SetLastBlockHeightByValidator(ctx, valAddr, proposal.GetHeight())

	return att, nil
}

func (k Keeper) TryAttestation(ctx sdk.Context, att *types.Attestation) {
	proposal, err := k.UnpackAttestationProposal(att)
	if err != nil {
		panic("could not cast to proposal")
	}
	hash, err := proposal.ProposalHash()
	if err != nil {
		panic("unable to compute proposal hash")
	}

	if !att.Observed {
		totalPower, err := k.StakingKeeper.GetLastTotalPower(ctx)
		if err != nil {
			panic(err)
		}
		requiredPower := types.AttestationVotesPowerThreshold.Mul(totalPower).Quo(math.NewInt(100))
		attestationPower := math.NewInt(0)

		proposalType := proposal.GetType()

		if proposalType == 1 || proposalType == 2 || proposalType == 3 {
			allValidators, err := k.StakingKeeper.GetAllValidators(ctx)
			if err != nil {
				panic(err)
			}
			activeValidatorCount := 0
			for _, validator := range allValidators {
				if validator.GetBondedTokens().IsPositive() {
					activeValidatorCount++
				}
			}

			receivedVotes := math.NewInt(int64(len(att.Votes)))
			votesNeeded := types.AttestationVoteCountThreshold.Mul(math.NewInt(int64(activeValidatorCount))).Quo(math.NewInt(100))

			if receivedVotes.GTE(votesNeeded) {
				att.Observed = true
				err := k.processAttestation(ctx, att, proposal)
				if err == nil {
					k.SetAttestation(ctx, proposal.GetHeight(), hash, att)
					k.emitObservedEvent(ctx, att, proposal)
				}
			}
		} else if proposalType == 0 {
			for _, validator := range att.Votes {
				val, err := sdk.ValAddressFromBech32(validator)
				if err != nil {
					panic(err)
				}
				validatorPower, err := k.StakingKeeper.GetLastValidatorPower(ctx, val)
				if err != nil {
					panic(err)
				}

				attestationPower = attestationPower.Add(math.NewInt(validatorPower))

				if attestationPower.GTE(requiredPower) {
					att.Observed = true
					err := k.processAttestation(ctx, att, proposal)
					if err == nil {
						k.SetAttestation(ctx, proposal.GetHeight(), hash, att)
						k.emitObservedEvent(ctx, att, proposal)
					}
					break
				}
			}
		}
	}
}

func (k Keeper) emitObservedEvent(ctx sdk.Context, att *types.Attestation, proposal types.BtcProposal) {
	hash, err := proposal.ProposalHash()
	if err != nil {
		panic(fmt.Errorf("unable to compute proposal hash: %w", err))
	}

	ctx.EventManager().EmitTypedEvent(
		&types.EventObservation{
			AttestationType: string(proposal.GetType()),
			AttestationId:   string(types.GetAttestationKey(proposal.GetHeight(), hash)),
		},
	)
}

func (k Keeper) SetAttestation(ctx sdk.Context, height uint64, proposalHash []byte, att *types.Attestation) {
	store := k.KVStore(ctx)
	aKey := types.GetAttestationKey(height, proposalHash)
	store.Set(aKey, k.cdc.MustMarshal(att))
}

func (k Keeper) GetAttestation(ctx sdk.Context, height uint64, proposalHash []byte) *types.Attestation {
	store := k.KVStore(ctx)
	aKey := types.GetAttestationKey(height, proposalHash)
	bz := store.Get(aKey)
	if len(bz) == 0 {
		return nil
	}
	var att types.Attestation
	k.cdc.MustUnmarshal(bz, &att)
	return &att
}

func (k Keeper) GetSweepProposalAttestationsForBtcSweepTx(ctx sdk.Context, txHash string) (types.Attestation, error) {
	var filteredAttestation types.Attestation
	found := false
	k.IterateAttestations(ctx, false, func(_ []byte, att types.Attestation) bool {
		proposal, err := k.UnpackAttestationProposal(&att)
		if err != nil {
			panic("couldn't cast to proposal")
		}
		if att.Observed && proposal.GetType() == types.PROPOSAL_TYPE_SWEEP_PROPOSAL {
			hash, err := proposal.ProposalHash()
			if err != nil {
				panic(fmt.Errorf("unable to compute proposal hash: %w", err))
			}
			txHashStr := hex.EncodeToString(hash)
			tx, err := types.CreateTxFromHex(txHashStr)
			if err != nil {
				panic(err)
			}
			txHashFromProposal := tx.TxHash().String()
			if txHash == txHashFromProposal {
				filteredAttestation = att
				found = true
				return true
			}
		}
		return false
	})

	if !found {
		return types.Attestation{}, fmt.Errorf("no matching attestation found for txHash: %s", txHash)
	}

	return filteredAttestation, nil
}

func (k Keeper) processAttestation(ctx sdk.Context, att *types.Attestation, proposal types.BtcProposal) error {
	hash, err := proposal.ProposalHash()
	if err != nil {
		panic(fmt.Errorf("unable to compute proposal hash: %w", err))
	}
	xCtx, commit := ctx.CacheContext()
	if err := k.AttestationHandler.Handle(xCtx, *att, proposal); err != nil {
		k.Logger(ctx).Error("attestation failed",
			"cause", err.Error(),
			"proposal type", proposal.GetType(),
			"id", types.GetAttestationKey(proposal.GetHeight(), hash),
		)
		return err
	} else {
		commit()
	}
	return nil
}

func (k Keeper) GetAttestationMapping(ctx sdk.Context) (attestationMapping map[uint64][]types.Attestation, orderedKeys []uint64) {
	attestationMapping = make(map[uint64][]types.Attestation)
	k.IterateAttestations(ctx, false, func(_ []byte, att types.Attestation) bool {
		proposal, err := k.UnpackAttestationProposal(&att)
		if err != nil {
			panic("couldn't cast to proposal")
		}

		if val, ok := attestationMapping[proposal.GetHeight()]; !ok {
			attestationMapping[proposal.GetHeight()] = []types.Attestation{att}
		} else {
			attestationMapping[proposal.GetHeight()] = append(val, att)
		}
		return false
	})
	orderedKeys = make([]uint64, 0, len(attestationMapping))
	for k := range attestationMapping {
		orderedKeys = append(orderedKeys, k)
	}
	sort.Slice(orderedKeys, func(i, j int) bool { return orderedKeys[i] < orderedKeys[j] })

	return
}

func (k Keeper) UnpackAttestationProposal(att *types.Attestation) (types.BtcProposal, error) {
	var msg types.BtcProposal
	err := k.cdc.UnpackAny(att.Proposal, &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (k Keeper) IterateAttestations(ctx sdk.Context, reverse bool, cb func([]byte, types.Attestation) bool) {
	store := k.KVStore(ctx)
	prefix := types.OracleAttestationKey

	var iter storetypes.Iterator
	if reverse {
		iter = store.ReverseIterator(prefixRange(prefix))
	} else {
		iter = store.Iterator(prefixRange(prefix))
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		att := types.Attestation{
			Observed: false,
			Votes:    []string{},
			Height:   0,
			Proposal: &codectypes.Any{},
		}
		k.cdc.MustUnmarshal(iter.Value(), &att)
		if cb(iter.Key(), att) {
			return
		}
	}
}

func (k Keeper) GetLastObservedBlockHeight(ctx sdk.Context) uint64 {
	store := k.KVStore(ctx)
	bytes := store.Get(types.LastObservedBlockHeightKey)
	if len(bytes) == 0 {
		return 0
	}
	return types.UInt64FromBytes(bytes)
}

func (k Keeper) SetLastBlockHeightByValidator(ctx sdk.Context, validator sdk.ValAddress, nonce uint64) {
	store := k.KVStore(ctx)
	store.Set(types.GetLastBlockHeightByValidatorKey(validator), types.UInt64Bytes(nonce))
}

func prefixRange(prefix []byte) ([]byte, []byte) {
	if prefix == nil {
		panic("nil key not allowed")
	}
	if len(prefix) == 0 {
		return nil, nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	l := len(end) - 1
	end[l]++
	for end[l] == 0 && l > 0 {
		l--
		end[l]++
	}
	if l == 0 && end[0] == 0 {
		end = nil
	}
	return prefix, end
}
