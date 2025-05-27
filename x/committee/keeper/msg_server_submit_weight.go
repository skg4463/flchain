package keeper

import (
	"context"

	"flchain/x/committee/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SubmitWeight(goCtx context.Context, msg *types.MsgSubmitWeight) (*types.MsgSubmitWeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgSubmitWeightResponse{}, nil
}
