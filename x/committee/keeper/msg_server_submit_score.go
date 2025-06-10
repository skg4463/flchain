package keeper

import (
	"context"

	"flchain/x/committee/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SubmitScore(goCtx context.Context, msg *types.MsgSubmitScore) (*types.MsgSubmitScoreResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgSubmitScoreResponse{}, nil
}
