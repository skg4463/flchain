package keeper

import (
	"flchain/x/flchain/types"
)

var _ types.QueryServer = Keeper{}
