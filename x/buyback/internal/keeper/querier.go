package keeper

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/e-money/em-ledger/x/buyback/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryBalance:
			return queryBalance(ctx, k)

		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unrecognized buyback query endpoint")
		}
	}
}

type QueryBalanceResponse struct {
	Balance sdk.Coins `json:"balance" yaml:"balance"`
}

func queryBalance(ctx sdk.Context, k Keeper) ([]byte, error) {
	account := k.GetBuybackAccount(ctx)

	response := QueryBalanceResponse{
		Balance: account.GetCoins(),
	}

	return json.Marshal(response)
}
