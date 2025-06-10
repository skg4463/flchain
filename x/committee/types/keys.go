package types

import (
	"encoding/binary"
)

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

// SubmissionKey는 round+lnodeId로 유니크하게 생성
func SubmissionKey(round uint64, lnodeId string) []byte {
	// round를 BigEndian uint64 바이트로 변환
	roundBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(roundBytes, round)
	// prefix + lnodeId + roundBytes 조합 (예: "submission:node1:\x00\x00...")
	return append([]byte("submission:"+lnodeId+":"), roundBytes...)
}
