package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMint{}

func NewMsgMint(requester string, amount uint64) *MsgMint {
	return &MsgMint{
		Requester: requester,
		Amount:    amount,
	}
}

func (msg *MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Requester)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid requester address (%s)", err)
	}

	// Validate amount
	if msg.Amount <= 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount must be positive, got %d", msg.Amount)
	}

	// You can add a maximum amount limit (adjust the value based on your requirements)
	//const maxAmount = uint64(1000000) // Example maximum amount
	//if msg.Amount > maxAmount {
	//	return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "amount exceeds maximum allowed (%d)", maxAmount)
	//}

	return nil

}
