// submit-weight 모듈
// 상태 저장 구조 정의
// L-node가 제출한 학습 결과를 저장하는 구조체 정의
// SPDX-License-Identifier: Apache-2.0
// 테스트

package types

import (
	"encoding/binary"
	//"fmt"
	//sdk "github.com/cosmos/cosmos-sdk/types"
)

// WeightSubmission L-node가 제출한 학습 결과 저장구조체
// creator == Lnodeid : 대리제출을 제외하고 원칙적으로 같을 것을 요구구
type WeightSubmission struct {
	Creator         string `json:"creator" yaml:"creator"`                   // tx sender
	LnodeId         string `json:"lnode_id" yaml:"lnode_id"`                 // L-node ID
	EncryptedWeight string `json:"encrypted_weight" yaml:"encrypted_weight"` // 암호화된 학습 결과
	Round           uint64 `json:"round" yaml:"round"`                       // 제출 라운드
}

// GetWeightSubmissionKey는 라운드 + L-nodeId 조합으로 고유한 키를 생성
func GetWeightSubmissionKey(round uint64, lnodeId string) []byte {
	return append(GetRoundPrefixKey(round), []byte(lnodeId)...)
}

// GetRoundPrefixKey는 라운드별 prefix key를 생성
// binary.BigEndian.PutUint64()을 쓰는 건 Cosmos SDK 키 정렬 기준과 맞추기 위함
func GetRoundPrefixKey(round uint64) []byte {
	prefix := []byte("weight:")
	roundBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(roundBytes, round)
	return append(prefix, roundBytes...)
}
