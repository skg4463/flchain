package keeper

import (
	"context"

	"flchain/x/committee/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetScore(goCtx context.Context, req *types.QueryGetScoreRequest) (*types.QueryGetScoreResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	key := types.ScoreKey(req.Round, req.CnodeId, req.LnodeId)
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(key)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get score from store")
	}
	if bz == nil {
		return nil, status.Error(codes.NotFound, "score not found")
	}

	// bz는 이미 JSON string이므로 그대로 반환
	return &types.QueryGetScoreResponse{
		ScoreJson: string(bz),
	}, nil
}
