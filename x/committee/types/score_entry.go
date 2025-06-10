package types

// ScoreEntryKey는 L-node ID를 기반으로 고유한 키를 생성
type ScoreEntry struct {
	LnodeId string `json:"lnodeId"`
	Score   string `json:"score"`
}
