package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "flchain/testutil/keeper"
	"flchain/x/flchain/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.FlchainKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
