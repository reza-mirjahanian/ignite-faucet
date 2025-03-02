package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"blog/x/faucet/types"
)

func (k Keeper) ListMinted(goCtx context.Context, req *types.QueryListMintedRequest) (*types.QueryListMintedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	//todo Paginate long result
	ctx := sdk.UnwrapSDKContext(goCtx)
	history, userTotalMinted := k.GetRequestHistory(ctx, sdk.AccAddress(req.Address))

	return &types.QueryListMintedResponse{Requests: &history, TotalAmount: userTotalMinted}, nil

}
