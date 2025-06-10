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

// ScoreKey는 라운드, cnodeId, lnodeId 조합으로 유니크한 key 생성
// cnodeId는 committee node ID, lnodeId는 learner node ID
// 예: "Score:Round:Cnode1:Lnode1"
func ScoreKey(round uint64, cnodeId, lnodeId string) []byte {
	roundBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(roundBytes, round)
	// 접두사("score:"), round, cnodeId, lnodeId 순으로 단순 이어붙임
	return append(append(append([]byte("score:"), roundBytes...), []byte(":"+cnodeId+":")...), []byte(lnodeId)...)
}
