package blog

import (
	"blog/x/faucet/keeper"
	"blog/x/faucet/types"
	"cosmossdk.io/math"
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdkTypes.Context, k keeper.Keeper, genState types.GenesisState) {
	fmt.Println(types.BrightCyan + "func InitGenesis() ")
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	/////////////////
	if acc := k.GetFaucetAccount(ctx); acc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	} else {
		fmt.Println(types.BrightYellow + "GetFaucetAccount: " + acc.String())
	}

	// Params maxPerRequest and maxPerAddress are modifiable via governance
	k.ActiveFaucet(ctx, true)
	k.SetSafeTimeout(ctx, 30*time.Second)
	k.SetFaucetMaxCapacity(ctx, math.NewInt(2000000000))
	k.SetMaxPerRequest(ctx, math.NewInt(20))
	k.SetMaxPerAddress(ctx, math.NewInt(1000))
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdkTypes.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	// this line is used by starport scaffolding # genesis/module/export
	return genesis
}

// ExportGenesis exports genesis state
//func ExportGenesis(ctx sdkTypes.Context, k Keeper) types.GenesisState {
//	return types.GenesisState{
//		EnableFaucet:        k.IsFaucetActive(ctx),
//		Timeout:             k.GetSafeTimeout(ctx),
//		FaucetCap:           k.GetFaucetMaxCapacity(ctx),
//		MaxAmountPerRequest: k.GetMaxPerRequest(ctx),
//	}
//}
