package keeper

import (
	"context"
	"encoding/json"

	"flchain/x/committee/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetCommitAtt(goCtx context.Context, req *types.QueryGetCommitAttRequest) (*types.QueryGetCommitAttResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	key := types.CommitAttKey(req.Round)
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(key)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get commit att state")
	}
	if bz == nil {
		return nil, status.Error(codes.NotFound, "commit att not found")
	}

	var commitAtt types.CommitAtt
	if err := json.Unmarshal(bz, &commitAtt); err != nil {
		return nil, status.Error(codes.Internal, "unmarshal error")
	}

	return &types.QueryGetCommitAttResponse{
		CommitAtt: &commitAtt,
	}, nil
}
