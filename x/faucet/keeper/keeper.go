package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
	"time"

	"blog/x/faucet/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority     string
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper

		timeouts map[string]time.Time
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,

) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetFaucetAccount returns the faucet ModuleAccount
//func (k Keeper) GetFaucetAccount(ctx sdk.Context) supplyexported.ModuleAccountI {
//	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
//}

func (k Keeper) GetFaucetAccount(ctx sdk.Context) sdk.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// Fund checks for timeout and max thresholds and then mints coins and transfers
// coins to the recipient.
func (k Keeper) Fund(ctx sdk.Context, amount sdk.Coins, recipient sdk.AccAddress) error {
	if !k.IsEnabled(ctx) {
		return errors.New("faucet is not enabled. Restart the application and set faucet's 'enable_faucet' genesis field to true")
	}

	if err := k.rateLimit(ctx, recipient.String()); err != nil {
		return err
	}

	totalRequested := sdkmath.ZeroInt()
	for _, coin := range amount {
		totalRequested = totalRequested.Add(coin.Amount)
	}

	maxPerReq := k.GetMaxPerRequest(ctx)
	if totalRequested.GT(maxPerReq) {
		return fmt.Errorf("canot fund more than %s per request. requested %s", maxPerReq, totalRequested)
	}

	funded := k.GetFunded(ctx)
	totalFunded := sdkmath.ZeroInt()
	for _, coin := range funded {
		totalFunded = totalFunded.Add(coin.Amount)
	}

	capacity := k.GetCap(ctx)

	if totalFunded.Add(totalRequested).GT(capacity) {
		return fmt.Errorf("maximum capacity of %s reached. Cannot continue funding", capacity)
	}

	//if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, amount); err != nil {
	//	return err
	//}
	//
	//if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, amount); err != nil {
	//	return err
	//}
	// Mint coins into the module account
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, amount); err != nil {
		return err
	}

	// Send minted coins from the module account to the recipient
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, recipient, amount); err != nil {
		return err
	}

	k.SetFunded(ctx, funded.Add(amount...))

	k.logger.Info(fmt.Sprintf("funded %s to %s", amount, recipient))
	return nil
}

func (k Keeper) GetTimeout(ctx sdk.Context) time.Duration {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.TimeoutKey)
	if bz == nil {
		return time.Duration(0)
	}
	return time.Duration(sdk.BigEndianToUint64(bz))
}

func (k Keeper) SetTimout(ctx sdk.Context, timeout time.Duration) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	modStore.Set(types.TimeoutKey, sdk.Uint64ToBigEndian(uint64(timeout)))
	//todo search to convert in Unit directly
}

func (k Keeper) IsEnabled(ctx sdk.Context) bool {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.EnableFaucetKey)
	if types.IsTrueB(bz) {
		return true
	} else {
		return false
	}
}

func (k Keeper) SetEnabled(ctx sdk.Context, enabled bool) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	val := types.ToBoolB(enabled)
	modStore.Set(types.EnableFaucetKey, []byte{val})
}

func (k Keeper) GetCap(ctx sdk.Context) sdkmath.Int {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.CapKey)
	if len(bz) == 0 {
		return sdkmath.ZeroInt()
	}
	var capValue big.Int
	return sdkmath.NewIntFromBigInt((&capValue).SetBytes(bz))
}

func (k Keeper) SetCap(ctx sdk.Context, cap sdkmath.Int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	modStore.Set(types.CapKey, cap.BigInt().Bytes())
}

func (k Keeper) GetMaxPerRequest(ctx sdk.Context) sdkmath.Int {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.MaxPerRequestKey)
	if len(bz) == 0 {
		return sdkmath.ZeroInt()
	} // todo Maybe better vale instead of zero!
	var maxPerReq big.Int
	return sdkmath.NewIntFromBigInt((&maxPerReq).SetBytes(bz))
}

func (k Keeper) SetMaxPerRequest(ctx sdk.Context, maxPerReq sdkmath.Int) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	modStore.Set(types.MaxPerRequestKey, maxPerReq.BigInt().Bytes())
}

func (k Keeper) GetFunded(ctx sdk.Context) sdk.Coins {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := modStore.Get(types.FundedKey)
	if len(bz) == 0 {
		return nil
	}

	var funded sdk.Coins
	k.cdc.MustUnmarshalBinaryBare(bz, &funded)

	return funded
}

func (k Keeper) SetFunded(ctx sdk.Context, funded sdk.Coins) {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	modStore := prefix.NewStore(storeAdapter, types.KeyPrefix(types.StoreKey))
	bz := k.cdc.MustMarshal(&funded)
	modStore.Set(types.FundedKey, bz)
}

func (k Keeper) rateLimit(ctx sdk.Context, address string) error {
	// first time requester, can send request
	lastRequest, ok := k.timeouts[address]
	if !ok {
		k.timeouts[address] = time.Now().UTC() // todo Maybe, use block time instead
		return nil
	}

	defaultTimeout := k.GetTimeout(ctx)
	sinceLastRequest := time.Since(lastRequest)

	if defaultTimeout > sinceLastRequest {
		wait := defaultTimeout - sinceLastRequest
		return fmt.Errorf("%s has requested funds within the last %s, wait %s before trying again", address, defaultTimeout.String(), wait.String())
	}

	// user able to send funds since they have waited for period
	k.timeouts[address] = time.Now().UTC()
	return nil
}
