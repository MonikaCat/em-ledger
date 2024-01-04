package keeper

import (
	"context"

	"github.com/MonikaCat/em-ledger/x/authority/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) GasPrices(c context.Context, req *types.QueryGasPricesRequest) (*types.QueryGasPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	gasPrices := k.GetGasPrices(ctx)
	return &types.QueryGasPricesResponse{MinGasPrices: gasPrices}, nil
}

func (k Keeper) UpgradePlan(c context.Context, req *types.QueryUpgradePlanRequest) (*types.QueryUpgradePlanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	plan, _ := k.GetUpgradePlan(ctx)

	return &types.QueryUpgradePlanResponse{Plan: plan}, nil
}
