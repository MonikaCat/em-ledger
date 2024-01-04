package keeper

import (
	"context"

	"github.com/MonikaCat/em-ledger/x/inflation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Inflation(c context.Context, req *types.QueryInflationRequest) (*types.QueryInflationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	response := types.QueryInflationResponse{
		State: k.GetState(ctx),
	}
	return &response, nil
}
