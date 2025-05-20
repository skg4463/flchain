package flchain_test

import (
	"testing"

	keepertest "flchain/testutil/keeper"
	"flchain/testutil/nullify"
	flchain "flchain/x/flchain/module"
	"flchain/x/flchain/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.FlchainKeeper(t)
	flchain.InitGenesis(ctx, k, genesisState)
	got := flchain.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
