package keeper

import (
	"context"

	"flchain/x/committee/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/cosmos/cosmos-sdk/types/errors" 오류수정
)

var (
	ErrCustomUnknownRequest = errorsmod.Register("committee", 1, "unknown request")
	ErrCustomConflict       = errorsmod.Register("committee", 2, "conflict error")
	ErrCustomTxDecode       = errorsmod.Register("committee", 3, "tx decode error")
)

func (k msgServer) SubmitWeight(
	goCtx context.Context,
	msg *types.MsgSubmitWeight,
) (*types.MsgSubmitWeightResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 1. 상태 저장용 Submission 구조체 생성
	submission := types.Submission{
		Creator:         msg.Creator,
		LnodeId:         msg.LnodeId,
		EncryptedWeight: msg.EncryptedWeight,
		Round:           msg.Round,
	}

	// 2. KVStoreService 기반의 store 열기 (v0.50+ 표준)
	//  Keeper의 storeService에 바로 접근
	//  storeService는 Keeper의 필드로 정의되어 있어야 함
	store := k.storeService.OpenKVStore(ctx)

	// 3. 저장할 key 생성 (round+lnodeId 조합)
	key := types.SubmissionKey(msg.Round, msg.LnodeId)

	// 4. 중복 제출 방지 (Has는 error를 반환하므로 에러처리 필요)
	has, err := store.Has(key)
	if err != nil {
		return nil, errorsmod.Wrapf(ErrCustomUnknownRequest, "failed to check key existence: %v", err)
	}
	if has {
		return nil, errorsmod.Wrapf(ErrCustomConflict, "already submitted by lnode %s in round %d", msg.LnodeId, msg.Round)
	}

	// 5. proto 구조체 marshal 및 저장
	bz, err := submission.Marshal()
	if err != nil {
		return nil, errorsmod.Wrapf(ErrCustomTxDecode, "cannot marshal submission: %v", err)
	}
	if err := store.Set(key, bz); err != nil {
		return nil, errorsmod.Wrapf(ErrCustomUnknownRequest, "failed to store submission: %v", err)
	}

	return &types.MsgSubmitWeightResponse{}, nil
}
