package types

const (
	// ModuleName defines the module name
	ModuleName = "flchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_flchain"
)

var (
	ParamsKey = []byte("p_flchain")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
