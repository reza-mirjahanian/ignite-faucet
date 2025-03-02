package types

import (
	"context"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper defines the expected interface for the Account module.
type AccountKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
	GetModuleAccount(ctx context.Context, moduleName string) sdk.ModuleAccountI
	NewAccount(context.Context, sdk.AccountI) sdk.AccountI
	SetAccount(goCtx context.Context, acc sdk.AccountI)
	SetModuleAccount(goCtx context.Context, acc sdk.ModuleAccountI)
	GetModuleAccountAndPermissions(ctx context.Context, moduleName string) (sdk.ModuleAccountI, []string)
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string)

	GetModulePermissions() map[string]types.PermissionsForAddress
	GetAllAccounts(ctx context.Context) []sdk.AccountI
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
	GetBalance(goCtx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(goCtx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(goCtx context.Context, name string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(goCtx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(goCtx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SetDenomMetaData(goCtx context.Context, denomMetaData banktypes.Metadata)
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
	SendCoins(goCtx context.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
