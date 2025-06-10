// submit-weight query

package keeper

import (
	"context"

	"flchain/x/committee/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetSubmission(goCtx context.Context, req *types.QueryGetSubmissionRequest) (*types.QueryGetSubmissionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// key 생성 (round, lnodeId)
	key := types.SubmissionKey(req.Round, req.LnodeId)

	// KVStoreService 기반 KVStore 열기
	store := k.storeService.OpenKVStore(ctx)

	// 값 조회
	bz, err := store.Get(key)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get submission")
	}
	if bz == nil {
		return nil, status.Error(codes.NotFound, "submission not found")
	}

	var submission types.Submission
	if err := submission.Unmarshal(bz); err != nil {
		return nil, status.Error(codes.Internal, "failed to unmarshal submission")
	}

	return &types.QueryGetSubmissionResponse{
		Submission: &submission,
	}, nil
}
