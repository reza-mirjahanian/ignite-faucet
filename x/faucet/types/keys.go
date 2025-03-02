package types

const (

	// ModuleName is the name of the module
	ModuleName = "faucet"

	// StoreKey to be used when creating the KVStore
	StoreKey              = ModuleName
	StoreHistoryRecordKey = "StoreHistoryRecordKey"

	//// RouterKey uses module name for tx routing
	//RouterKey = ModuleName
	//
	//// QuerierRoute uses module name for query routing
	//QuerierRoute = ModuleName
)

var (
	EnableFaucetKey  = []byte{0x01}
	TimeoutKey       = []byte{0x02}
	CapKey           = []byte{0x03}
	MaxPerRequestKey = []byte{0x04}
	TotalMintedKey   = []byte{0x05}
	MaxPerAddressKey = []byte{0x06}

	ParamsKey = []byte("p_faucet")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	// TrueB is a byte with value 1 that represents true.
	TrueB = byte(0x01)
	// FalseB is a byte with value 0 that represents false.
	FalseB = byte(0x00)

	FaucetDenom = "stake"
)

// IsTrueB returns true if the provided byte slice has exactly one byte, and it is equal to TrueB.
func IsTrueB(bz []byte) bool {
	return len(bz) == 1 && bz[0] == TrueB
}

// ToBoolB returns TrueB if v is true, and FalseB if it's false.
func ToBoolB(v bool) byte {
	if v {
		return TrueB
	}
	return FalseB
}
