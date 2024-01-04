package keeper

import (
	"testing"

	"github.com/MonikaCat/em-ledger/x/issuer/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryIssuers(t *testing.T) {
	encConfig := MakeTestEncodingConfig()
	ctx, _, _, keeper, _ := createTestComponentsWithEncodingConfig(t, encConfig)
	myIssuers := []types.Issuer{types.NewIssuer(sdk.AccAddress("emoney1n5ggspeff4fxc87dvmg0ematr3qzw5l4v20mdv"), "foo", "bar")}
	keeper.setIssuers(ctx, myIssuers)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encConfig.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, keeper)
	queryClient := types.NewQueryClient(queryHelper)

	specs := map[string]struct {
		req        *types.QueryIssuersRequest
		expIssuers []types.Issuer
	}{
		"all good": {
			req:        &types.QueryIssuersRequest{},
			expIssuers: myIssuers,
		},
		"nil param": {
			expIssuers: myIssuers,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			gotRsp, gotErr := queryClient.Issuers(sdk.WrapSDKContext(ctx), spec.req)
			require.NoError(t, gotErr)
			require.NotNil(t, gotRsp)
			assert.Equal(t, spec.expIssuers, gotRsp.Issuers)
		})
	}
}
