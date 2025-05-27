package keeper

import (
	"flchain/x/committee/types"
)

var _ types.QueryServer = Keeper{}
