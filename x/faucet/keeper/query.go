package keeper

import (
	"blog/x/faucet/types"
)

var _ types.QueryServer = Keeper{}
