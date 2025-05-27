package types

const (
	// ModuleName defines the module name
	ModuleName = "committee"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_committee"
)

var (
	ParamsKey = []byte("p_committee")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
