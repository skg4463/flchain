package committee_test

import (
	"testing"

	keepertest "flchain/testutil/keeper"
	"flchain/testutil/nullify"
	committee "flchain/x/committee/module"
	"flchain/x/committee/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.CommitteeKeeper(t)
	committee.InitGenesis(ctx, k, genesisState)
	got := committee.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
