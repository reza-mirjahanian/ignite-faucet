package keeper

import (
	"context"
	"cosmossdk.io/math"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"blog/x/faucet/types"
)

func (k msgServer) Mint(goCtx context.Context, msg *types.MsgMint) (*types.MsgMintResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	blockHeight := ctx.BlockHeight()
	requester := msg.Requester
	amount := msg.Amount
	fmt.Println(types.BrightMagenta+"Requester address is: "+types.Reset, requester)
	fmt.Println(types.BrightMagenta+"Request amount is: "+types.Reset, amount)
	fmt.Println(types.BrightMagenta+"BlockHeight is: "+types.Reset, blockHeight)

	coinAmount := sdk.Coin{
		Denom:  types.FaucetDenom,
		Amount: math.NewIntFromUint64(amount),
	}

	err := k.MintAndTransfer(
		ctx,
		coinAmount,
		sdk.AccAddress(requester),
	)
	if err != nil {
		fmt.Println(types.BrightRed+"Error in Mint: "+types.Reset, err)
		return nil, err
	}
	return &types.MsgMintResponse{
		Amount: amount,
	}, nil
}
