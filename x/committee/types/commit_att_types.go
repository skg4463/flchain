package types

// CommitAtt: 집계 결과(EndBlocker에서 계산)
type CommitAtt struct {
	Round     uint64             `json:"round"`     // 라운드 번호
	EwmaMap   map[string]float64 `json:"ewmaMap"`   // l-node별 EWMA
	SwMap     map[string]float64 `json:"swMap"`     // l-node별 합성가중치
	Ranking   []string           `json:"ranking"`   // Sw 내림차순 l-node 순위
	ClNode    string             `json:"clNode"`    // CL-node (1위)
	Committee []string           `json:"committee"` // 위원회(상위 5명)
}

//
// ScoreEntry: C-node의 평가 결과
// type ScoreEntry struct {
// 	LnodeId string `json:"lnodeId"`
// 	Score   string `json:"score"` // string 타입(파싱해서 float 사용)
// }
