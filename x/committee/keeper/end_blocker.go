package keeper

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"flchain/x/committee/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker: 블록 종료 시 집계, 위원회/CL-node 선출, 상태 저장
func (k Keeper) EndBlocker(ctx sdk.Context) {
	blockHeight := ctx.BlockHeight()

	// 라운드-블록 매핑: 3블록=1라운드 예시(환경에 맞게 조정)
	round := uint64(blockHeight / 3)

	// Block 역할 구분
	// Block 1(3n+1): l-node 제출 블록 (EndBlocker에서는 별도 작업 없음)
	// Block 2(3n+2): c-node 제출/집계 블록 (여기서 집계, 위원회/CL-node 선출)
	// Block 3(3n+3): validator set 교체 위한 staking tx (EndBlocker는 별도 작업 없음)

	// 역할 설명
	blockPhase := blockHeight%3 + 1
	var phaseDesc string
	switch blockPhase {
	case 1:
		phaseDesc = "L-node submit Phase"
	case 2:
		phaseDesc = "C-node submit/Aggregation Phase"
	case 3:
		phaseDesc = "Validator Set Replacement/Waiting Phase"
	}

	// 로그: 블록, 라운드, 라운드 역할 출력
	ctx.Logger().Info("test Round Config",
		"Round", round,
		"BlockPhase", blockPhase,
		"Role", phaseDesc,
		"Block", blockHeight,
	)

	// Logging: 라운드 정보 출력
	//fmt.Println("[Committee EndBlocker] Called for round", round)
	//ctx.Logger().Error("EndBlocker called for round", "round", round, "BlockHeight", blockHeight)

	if blockPhase == 2 {

		// 1. 이번 라운드의 l-node 목록 동적으로 불러오기
		lnodeIds := k.GetAllLnodeIds(ctx, round)

		// 2. 해당 라운드의 모든 l-node score 불러오기
		scores := k.GetAllScoresForRound(ctx, round)

		// 3. 이전 라운드 EWMA 값 불러오기 (없으면 0)
		prevEwma := k.GetPrevEwmaForAllLnodes(ctx, round-1, lnodeIds)

		// 4. EWMA 계산 및 저장
		beta := 0.7
		ewmaMap := map[string]float64{}
		sumEwma := 0.0
		for _, lnodeId := range lnodeIds {
			score := scores[lnodeId]
			oldEwma := prevEwma[lnodeId]
			ewma := beta*score + (1-beta)*oldEwma
			ewmaMap[lnodeId] = ewma
			sumEwma += ewma
			bz, _ := json.Marshal(ewma)
			k.Set(ctx, types.EwmaKey(round, lnodeId), bz)
		}

		// 5. Sw 계산 및 저장
		swMap := map[string]float64{}
		for _, lnodeId := range lnodeIds {
			ewma := ewmaMap[lnodeId]
			sw := 0.0
			if sumEwma > 0 {
				sw = ewma / sumEwma
			}
			swMap[lnodeId] = sw
			bz, _ := json.Marshal(sw)
			k.Set(ctx, types.SwKey(round, lnodeId), bz)
		}

		// 6. Sw 내림차순으로 l-node 랭킹 산출
		ranking := SortLnodesBySwDesc(swMap)

		// 7. committee/CL-node 선출
		clNode := ""
		committee := []string{}
		if len(ranking) > 0 {
			clNode = ranking[0]
			top := 5
			if len(ranking) < top {
				top = len(ranking)
			}
			committee = ranking[:top]
		}

		// 8. 집계 결과(CommitAtt)를 state에 JSON 직렬화로 저장
		k.StoreCommitAtt(ctx, round, ewmaMap, swMap, ranking, clNode, committee)

		// 9. 이벤트/로그 기록
		ctx.EventManager().EmitEvent(
			sdk.NewEvent("EndBlock-Committee",
				sdk.NewAttribute("round", fmt.Sprint(round)),
				sdk.NewAttribute("cl-node", clNode),
				sdk.NewAttribute("committee", fmt.Sprint(committee)),
			),
		)

		ctx.Logger().Info("test Aggregation Result",
			"round", round,
			"phase", phaseDesc,
			"phase", phaseDesc,
			"block", blockHeight,
			"cl-node(블록생성자)", clNode,
			"committee", committee,
			"ranking", ranking,
			"ewma", ewmaMap,
			"sw", swMap,
		)
	}
}

// ------------------- 유틸리티 함수들 -------------------

// Set: state에 키-값 저장 (KVStore 추상화)
func (k Keeper) Set(ctx sdk.Context, key, value []byte) {
	store := k.storeService.OpenKVStore(ctx)
	_ = store.Set(key, value)
}

// submit-weight 처리 시 호출: round별 l-node 목록에 추가
func (k Keeper) SetLnodeForRound(ctx sdk.Context, round uint64, lnodeId string) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.LnodeSetKey(round)
	var lnodes map[string]bool
	bz, _ := store.Get(key)
	if bz != nil {
		_ = json.Unmarshal(bz, &lnodes)
	} else {
		lnodes = map[string]bool{}
	}
	lnodes[lnodeId] = true
	bz, _ = json.Marshal(lnodes)
	k.Set(ctx, key, bz)
}

// round별 l-node id 동적 목록 반환
func (k Keeper) GetAllLnodeIds(ctx sdk.Context, round uint64) []string {
	store := k.storeService.OpenKVStore(ctx)
	key := types.LnodeSetKey(round)
	bz, _ := store.Get(key)
	lnodes := map[string]bool{}
	if bz != nil {
		_ = json.Unmarshal(bz, &lnodes)
	}
	result := []string{}
	for id := range lnodes {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}

// round별 모든 l-node score 가져오기
func (k Keeper) GetAllScoresForRound(ctx sdk.Context, round uint64) map[string]float64 {
	store := k.storeService.OpenKVStore(ctx)
	scores := make(map[string]float64)
	prefix := []byte(fmt.Sprintf("score:%d:", round))
	iter, _ := store.Iterator(prefix, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var scoreEntry types.ScoreEntry
		_ = json.Unmarshal(iter.Value(), &scoreEntry)
		scoreVal, _ := parseScoreString(scoreEntry.Score)
		lnodeId := scoreEntry.LnodeId
		scores[lnodeId] = scoreVal
	}
	return scores
}

// 이전 라운드의 모든 l-node ewma 값 가져오기
func (k Keeper) GetPrevEwmaForAllLnodes(ctx sdk.Context, prevRound uint64, lnodeIds []string) map[string]float64 {
	store := k.storeService.OpenKVStore(ctx)
	ewmas := make(map[string]float64)
	for _, lnodeId := range lnodeIds {
		key := types.EwmaKey(prevRound, lnodeId)
		bz, _ := store.Get(key)
		val := 0.0
		if bz != nil {
			_ = json.Unmarshal(bz, &val)
		}
		ewmas[lnodeId] = val
	}
	return ewmas
}

// Sw 내림차순 l-node 정렬
func SortLnodesBySwDesc(swMap map[string]float64) []string {
	type kv struct {
		Key   string
		Value float64
	}
	var sorted []kv
	for k, v := range swMap {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Value > sorted[j].Value })
	result := []string{}
	for _, kv := range sorted {
		result = append(result, kv.Key)
	}
	return result
}

// 집계 결과(CommitAtt)를 state에 JSON 직렬화로 저장
func (k Keeper) StoreCommitAtt(ctx sdk.Context, round uint64, ewmaMap, swMap map[string]float64, ranking []string, clNode string, committee []string) {
	//store := k.storeService.OpenKVStore(ctx)
	result := types.CommitAtt{
		Round:     round,
		EwmaMap:   ewmaMap,
		SwMap:     swMap,
		Ranking:   ranking,
		ClNode:    clNode,
		Committee: committee,
	}
	bz, _ := json.Marshal(result)
	k.Set(ctx, types.CommitAttKey(round), bz)
	k.Set(ctx, types.CommitteeKey(round), []byte(clNode))
}

// string 타입 점수 파싱
func parseScoreString(scoreStr string) (float64, error) {
	return strconv.ParseFloat(scoreStr, 64)
}
