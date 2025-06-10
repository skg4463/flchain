package keeper

import (
	"context"
	"encoding/json"

	"flchain/x/committee/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ErrCustomUnknownRequest1 = errorsmod.Register("committee", 11, "unknown request in submit-score")
	ErrCustomScoreDecode     = errorsmod.Register("committee", 12, "score decode error in submit-score")
	ErrCustomScoreConflict   = errorsmod.Register("committee", 13, "score conflict error in submit-score")
)

func (k msgServer) SubmitScore(
	goCtx context.Context,
	msg *types.MsgSubmitScore,
) (*types.MsgSubmitScoreResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. scoresJson을 Go 구조체 슬라이스로 언마샬
	var scores []types.ScoreEntry
	if err := json.Unmarshal([]byte(msg.ScoresJson), &scores); err != nil {
		return nil, errorsmod.Wrapf(ErrCustomScoreDecode, "invalid scoresJson: %v", err)
	}

	store := k.storeService.OpenKVStore(ctx)

	// 2. 각 ScoreEntry별로 저장(중복 방지)
	for _, score := range scores {
		key := types.ScoreKey(msg.Round, msg.CnodeId, score.LnodeId)
		has, err := store.Has(key)
		if err != nil {
			return nil, errorsmod.Wrapf(ErrCustomUnknownRequest, "failed to check key existence: %v", err)
		}
		if has {
			return nil, errorsmod.Wrapf(ErrCustomScoreConflict, "score for cnode %s, lnode %s, round %d already exists", msg.CnodeId, score.LnodeId, msg.Round)
		}
		bz, err := json.Marshal(score)
		if err != nil {
			return nil, errorsmod.Wrapf(ErrCustomScoreDecode, "marshal error: %v", err)
		}
		if err := store.Set(key, bz); err != nil {
			return nil, errorsmod.Wrapf(ErrCustomUnknownRequest, "failed to store score: %v", err)
		}
	}

	return &types.MsgSubmitScoreResponse{}, nil
}
