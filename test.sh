#!/bin/bash

CHAIN_ID="flchain"
KEYRING="test"
FEE="5000stake"
L_NODES=("alice" "bob" "c")
C_NODES=("alice" "bob" "c")

while true; do
    BLOCK_HEIGHT=$(curl -s localhost:26657/status | jq -r .result.sync_info.latest_block_height)
    PHASE=$((BLOCK_HEIGHT % 3))
    ROUND=$((BLOCK_HEIGHT / 3))

    # 페이즈1: l-node 제출
    if [ $PHASE -eq 1 ]; then
        echo "[Phase1] Block $BLOCK_HEIGHT: l-node 제출 (ROUND $ROUND)"
        for LNODE in "${L_NODES[@]}"; do
            flchaind tx committee submit-weight $LNODE "weight-$RANDOM" $ROUND \
                --from $LNODE --keyring-backend $KEYRING --chain-id $CHAIN_ID --fees $FEE --yes
        done
    fi

    # 페이즈2: c-node 평가 제출 (랜덤 점수)
    if [ $PHASE -eq 2 ]; then
        echo "[Phase2] Block $BLOCK_HEIGHT: c-node 평가/제출 (ROUND $ROUND)"
        for CNODE in "${C_NODES[@]}"; do
            SCORE_A=$(awk "BEGIN { printf \"%.2f\", $(($RANDOM%80+20))/100 }")
            SCORE_B=$(awk "BEGIN { printf \"%.2f\", $(($RANDOM%80+20))/100 }")
            SCORE_C=$(awk "BEGIN { printf \"%.2f\", $(($RANDOM%80+20))/100 }")
            SCORES_JSON="[\
{\"lnodeId\":\"alice\",\"score\":\"$SCORE_A\"},\
{\"lnodeId\":\"bob\",\"score\":\"$SCORE_B\"},\
{\"lnodeId\":\"c\",\"score\":\"$SCORE_C\"}]"
            flchaind tx committee submit-score $CNODE $ROUND "$SCORES_JSON" \
                --from $CNODE --keyring-backend $KEYRING --chain-id $CHAIN_ID --fees $FEE --yes
        done
    fi

    # 페이즈3: validator 변경 (CL-node만 남기기)
    if [ $PHASE -eq 0 ] && [ $BLOCK_HEIGHT -ne 0 ]; then
        echo "[Phase3] Block $BLOCK_HEIGHT: validator set 교체 (ROUND $ROUND)"
        # 집계 결과 쿼리로 cl-node 주소 추출
        CLNODE=$(flchaind query committee get-commit-att $ROUND | jq -r '.commitAtt.clNode')
        # 현재 validator 목록 추출
        VALS=$(flchaind query staking validators --output json | jq -r '.validators[].operator_address')
        for VAL in $VALS; do
            if [[ "$VAL" != "$CLNODE" ]]; then
                # cl-node 이외 모든 validator unbond (수량은 실험상 필요만큼 조절)
                flchaind tx staking unbond $VAL "1000000stake" --from $VAL --keyring-backend $KEYRING --chain-id $CHAIN_ID --yes --fees $FEE
            fi
        done
        # cl-node가 validator 아니면(처음 라운드 등), 등록
        # (이미 validator라면 오류 없이 넘어감)
        # pubkey 구하는 방법은 체인마다 다를 수 있으니 실험환경 맞춰 보완 필요
        # 예시: 
        # flchaind tx staking create-validator --amount="1000000stake" --pubkey=$(flchaind tendermint show-validator) --moniker="clnode" --chain-id=$CHAIN_ID --from=$CLNODE --keyring-backend=$KEYRING --fees=$FEE --yes
    fi

    sleep 5 # block 생성 주기 맞춰 조절
done
