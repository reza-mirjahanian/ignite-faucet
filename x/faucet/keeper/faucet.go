package keeper

import (
	"blog/x/faucet/types"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"time"
)

func (k Keeper) GetFaucetAccount(ctx sdk.Context) sdk.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// MintAndTransfer Main function
func (k Keeper) MintAndTransfer(ctx sdk.Context, coin sdk.Coin, recipient sdk.AccAddress) error {
	fmt.Println(types.BrightYellow+"SendCoinsFromModuleToAccount() started..."+types.Reset, recipient, k.bankKeeper.GetAllBalances(ctx, recipient))
	if !k.IsFaucetActive(ctx) {
		return errors.New("faucet is not enabled. Restart the application and set faucet's 'enable_faucet' genesis field to true")
	}
	// Check  max thresholds (safety)
	if err := k.rateLimitChecker(ctx, recipient.String()); err != nil {
		fmt.Println(types.BrightRed+"Error in rateLimitChecker "+types.Reset, err)
		return err
	}

	userRequestAmount := sdkmath.ZeroInt().Add(coin.Amount)

	// Check max thresholds (safety)
	maxPerReq := k.GetMaxPerRequest(ctx)
	if userRequestAmount.GT(maxPerReq) {
		return fmt.Errorf("canot fund more than %s per request. requested %s", maxPerReq, userRequestAmount)
	}

	//Check maxPerAddress (safety)
	history, userTotalMinted := k.GetRequestHistory(ctx, recipient)
	fmt.Println(types.BrightBlue+"User History, Total: "+types.Reset, history, userTotalMinted)
	maxUserMint := k.GetMaxPerAddress(ctx)
	if sdkmath.NewInt(int64(userTotalMinted)).Add(userRequestAmount).GT(maxUserMint) {
		return fmt.Errorf("canot fund more than %s for a user. User total minted until now: %d ", maxUserMint, userTotalMinted)
	}

	mintedCoins := k.GetFaucetTotalMinted(ctx)
	totalFaucetMint := sdkmath.ZeroInt().Add(mintedCoins.Amount)
	faucetMaxCapacity := k.GetFaucetMaxCapacity(ctx)
	if totalFaucetMint.Add(userRequestAmount).GT(faucetMaxCapacity) {
		return fmt.Errorf("maximum capacity of %s reached. Cannot continue funding", faucetMaxCapacity)
	}

	// Mint coins into the module account
	coins := sdk.Coins{
		coin,
	}
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}
	fmt.Println(types.BrightYellow+"After MintCoins(): "+types.Reset, coins)
	// Send mintedCoins coins from the module account to the recipient
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, coins); err != nil {
		return err
	}

	k.SetFaucetTotalMinted(ctx, mintedCoins.Add(coin))
	request := types.Request{Height: uint64(ctx.BlockHeight()), Amount: userRequestAmount.Uint64(), CreatedAt: uint64(ctx.BlockTime().UTC().Unix())}
	k.AppendRequestHistory(ctx, request, recipient)
	fmt.Println(types.BrightYellow+"SendCoinsFromModuleToAccount() done: "+types.Reset, recipient, k.bankKeeper.GetAllBalances(ctx, recipient))
	k.logger.Info(fmt.Sprintf("mintedCoins %s to %s", coin, recipient))
	return nil
}

func (k Keeper) GetSafeTimeout(ctx sdk.Context) time.Duration {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.TimeoutKey)
	if bz == nil {
		return time.Duration(0)
	}
	return time.Duration(sdk.BigEndianToUint64(bz))
}

func (k Keeper) SetSafeTimeout(ctx sdk.Context, timeout time.Duration) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	modStore.Set(types.TimeoutKey, sdk.Uint64ToBigEndian(uint64(timeout)))
}

func (k Keeper) IsFaucetActive(ctx sdk.Context) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.EnableFaucetKey)
	if types.IsTrueB(bz) {
		return true
	} else {
		return false
	}
}

func (k Keeper) ActiveFaucet(ctx sdk.Context, enabled bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	val := types.ToBoolB(enabled)
	modStore.Set(types.EnableFaucetKey, []byte{val})
}

func (k Keeper) GetFaucetMaxCapacity(ctx sdk.Context) sdkmath.Int {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.CapKey)
	if len(bz) == 0 {
		return sdkmath.ZeroInt()
	}
	var capValue big.Int
	return sdkmath.NewIntFromBigInt((&capValue).SetBytes(bz))
}

func (k Keeper) SetFaucetMaxCapacity(ctx sdk.Context, cap sdkmath.Int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	modStore.Set(types.CapKey, cap.BigInt().Bytes())
}

// Params maxPerRequest
func (k Keeper) GetMaxPerRequest(ctx sdk.Context) sdkmath.Int {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.MaxPerRequestKey)
	if len(bz) == 0 {
		return sdkmath.ZeroInt()
	} // todo Maybe better vale instead of zero! like default
	var maxPerReq big.Int
	return sdkmath.NewIntFromBigInt((&maxPerReq).SetBytes(bz))
}

func (k Keeper) SetMaxPerRequest(ctx sdk.Context, maxPerReq sdkmath.Int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	modStore.Set(types.MaxPerRequestKey, maxPerReq.BigInt().Bytes())
}

// Params maxPerAddress, refers to the maximum in total requested by the user in all time
func (k Keeper) GetMaxPerAddress(ctx sdk.Context) sdkmath.Int {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.MaxPerAddressKey)
	if len(bz) == 0 {
		return sdkmath.ZeroInt()
	} // todo Maybe better vale instead of zero! like default
	var maxPerReq big.Int
	return sdkmath.NewIntFromBigInt((&maxPerReq).SetBytes(bz))
}

func (k Keeper) SetMaxPerAddress(ctx sdk.Context, maxPerReq sdkmath.Int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	modStore.Set(types.MaxPerAddressKey, maxPerReq.BigInt().Bytes())
}

func (k Keeper) GetFaucetTotalMinted(ctx sdk.Context) sdk.Coin {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.TotalMintedKey)
	if len(bz) == 0 {
		return sdk.Coin{
			Denom:  types.FaucetDenom,
			Amount: sdkmath.NewInt(0),
		}
	}
	var funded sdk.Coin
	k.cdc.MustUnmarshal(bz, &funded)
	return funded
}

func (k Keeper) SetFaucetTotalMinted(ctx sdk.Context, funded sdk.Coin) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := k.cdc.MustMarshal(&funded)
	modStore.Set(types.TotalMintedKey, bz)
}

// Safety check
func (k Keeper) rateLimitChecker(ctx sdk.Context, address string) error {
	lastRequest, ok := k.timeoutsTable[address]
	if !ok {
		// For the first time
		k.timeoutsTable[address] = time.Now().UTC() // todo Maybe, use block time instead
		fmt.Println(types.BrightYellow+"k.timeoutsTable:  "+types.Reset, k.timeoutsTable)
		return nil
	}

	safeTimeout := k.GetSafeTimeout(ctx)
	sinceLastRequest := time.Since(lastRequest)

	if safeTimeout > sinceLastRequest {
		passed := safeTimeout - sinceLastRequest
		return fmt.Errorf("%s has requested funds within the last %s, wait %s before trying again", address, safeTimeout.String(), passed.String())
	}

	// user able to send funds since they have waited for period
	k.timeoutsTable[address] = time.Now().UTC()
	fmt.Println(types.BrightYellow+"k.timeoutsTable:  "+types.Reset, k.timeoutsTable)
	return nil
}

func (k Keeper) AppendRequestHistory(ctx sdk.Context, request types.Request, recipient sdk.AccAddress) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreHistoryRecordKey))
	existingHistory, _ := k.GetRequestHistory(ctx, recipient)
	// Create or update requests
	var requests types.Requests
	if len(existingHistory.Items) == 0 {
		requests = types.Requests{
			Items: []*types.Request{&request},
		}
	} else {
		requests = types.Requests{
			Items: append(existingHistory.Items, &request),
		}
	}
	appendedValue := k.cdc.MustMarshal(&requests)
	store.Set(recipient, appendedValue)
}

func (k Keeper) GetRequestHistory(ctx sdk.Context, recipient sdk.AccAddress) (types.Requests, uint64) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreHistoryRecordKey))
	bz := store.Get(recipient)
	if bz == nil {
		return types.Requests{}, 0
	}

	var userHistory types.Requests
	k.cdc.MustUnmarshal(bz, &userHistory)

	// Calculate total amount
	var totalAmount uint64
	for _, request := range userHistory.Items {
		totalAmount += request.Amount
	}

	return userHistory, totalAmount
}
