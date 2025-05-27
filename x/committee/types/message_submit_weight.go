package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSubmitWeight{}

func NewMsgSubmitWeight(creator string, lnodeId string, encryptedWeight string, round uint64) *MsgSubmitWeight {
	return &MsgSubmitWeight{
		Creator:         creator,
		LnodeId:         lnodeId,
		EncryptedWeight: encryptedWeight,
		Round:           round,
	}
}

func (msg *MsgSubmitWeight) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
