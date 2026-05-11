package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"twilight-project/nyks/x/forks/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const QUERY_ATTESTATIONS_LIMIT uint64 = 1000

func (q queryServer) GetAttestations(goCtx context.Context, req *types.QueryAttestationsRequest) (*types.QueryAttestationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	limit := req.Limit
	if limit == 0 || limit > QUERY_ATTESTATIONS_LIMIT {
		limit = QUERY_ATTESTATIONS_LIMIT
	}

	var (
		attestations []types.Attestation
		count        uint64
		iterErr      error
	)

	reverse := strings.EqualFold(req.OrderBy, "desc")
	filter := req.BtcHeight > 0 || req.ProposalType != ""

	q.k.IterateAttestations(ctx, reverse, func(_ []byte, att types.Attestation) (abort bool) {
		proposal, err := q.k.UnpackAttestationProposal(&att)
		if err != nil {
			iterErr = err
			return true
		}

		var match bool
		switch {
		case filter && proposal.GetHeight() == req.BtcHeight:
			attestations = append(attestations, att)
			match = true
		case filter && proposal.GetType().String() == req.ProposalType:
			attestations = append(attestations, att)
			match = true
		case !filter:
			attestations = append(attestations, att)
			match = true
		}

		if match {
			count++
			if count >= limit {
				return true
			}
		}

		return false
	})
	if iterErr != nil {
		return nil, iterErr
	}

	return &types.QueryAttestationsResponse{Attestations: attestations}, nil
}
