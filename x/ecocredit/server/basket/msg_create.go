package basket

import (
	"context"

	"github.com/cosmos/cosmos-sdk/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	api "github.com/regen-network/regen-ledger/api/regen/ecocredit/basket/v1"
	"github.com/regen-network/regen-ledger/x/ecocredit"
	"github.com/regen-network/regen-ledger/x/ecocredit/basket"
	"github.com/regen-network/regen-ledger/x/ecocredit/core"
)

// Create is an RPC to handle basket.MsgCreate
func (k Keeper) Create(ctx context.Context, msg *basket.MsgCreate) (*basket.MsgCreateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	var fee sdk.Coins
	k.paramsKeeper.Get(sdkCtx, core.KeyBasketCreationFee, &fee)
	if !msg.Fee.IsAllGTE(fee) {
		return nil, sdkerrors.ErrInsufficientFee.Wrapf("minimum fee %s, got %s", fee, msg.Fee)
	}

	curator, err := sdk.AccAddressFromBech32(msg.Curator)
	if err != nil {
		return nil, err
	}

	for _, coin := range fee {
		curatorBalance := k.bankKeeper.GetBalance(sdkCtx, curator, coin.Denom)
		if curatorBalance.IsNil() || curatorBalance.IsLT(coin) {
			return nil, sdkerrors.ErrInsufficientFunds.Wrapf("insufficient balance for bank denom %s", coin.Denom)
		}
	}

	err = k.distKeeper.FundCommunityPool(sdkCtx, fee, curator)
	if err != nil {
		return nil, err
	}

	creditType, err := k.coreStore.CreditTypeTable().Get(ctx, msg.CreditTypeAbbrev)
	if err != nil {
		return nil, sdkerrors.ErrInvalidRequest.Wrapf(
			"could not get credit type with abbreviation %s: %s", msg.CreditTypeAbbrev, err.Error(),
		)
	}

	denom, displayDenom, err := basket.FormatBasketDenom(msg.Name, msg.CreditTypeAbbrev, creditType.Precision)
	if err != nil {
		return nil, err
	}

	id, err := k.stateStore.BasketTable().InsertReturningID(ctx, &api.Basket{
		Curator:           curator,
		BasketDenom:       denom,
		DisableAutoRetire: msg.DisableAutoRetire,
		CreditTypeAbbrev:  msg.CreditTypeAbbrev,
		DateCriteria:      msg.DateCriteria.ToApi(),
		Exponent:          creditType.Precision, // exponent is no longer used but set until removed
		Name:              msg.Name,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "basket with name %s already exists", msg.Name)
	}
	if err = k.indexAllowedClasses(ctx, id, msg.AllowedClasses, msg.CreditTypeAbbrev); err != nil {
		return nil, err
	}

	denomUnits := []*banktypes.DenomUnit{{
		Denom:    displayDenom,
		Exponent: creditType.Precision,
	}}
	if creditType.Precision != 0 {
		denomUnits = append(denomUnits, &banktypes.DenomUnit{
			Denom: denom,
		})
	}

	k.bankKeeper.SetDenomMetaData(sdkCtx, banktypes.Metadata{
		DenomUnits:  denomUnits,
		Description: msg.Description,
		Base:        denom,
		Display:     displayDenom,
		Name:        msg.Name,
		Symbol:      msg.Name,
	})

	err = sdkCtx.EventManager().EmitTypedEvent(&basket.EventCreate{
		BasketDenom: denom,
		Curator:     msg.Curator,
	})

	return &basket.MsgCreateResponse{BasketDenom: denom}, err
}

// indexAllowedClasses checks that all `allowedClasses` both exist, and are of the specified credit type, then inserts
// the class into the BasketClass table.
func (k Keeper) indexAllowedClasses(ctx context.Context, basketID uint64, allowedClasses []string, creditTypeAbbrev string) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	for _, class := range allowedClasses {
		classInfo, err := k.coreStore.ClassTable().GetById(ctx, class)
		if err != nil {
			return sdkerrors.ErrInvalidRequest.Wrapf("could not get credit class %s: %s", class, err.Error())
		}

		if classInfo.CreditTypeAbbrev != creditTypeAbbrev {
			return sdkerrors.ErrInvalidRequest.Wrapf("basket specified credit type %s, but class %s is of type %s",
				creditTypeAbbrev, class, classInfo.CreditTypeAbbrev)
		}

		if err := k.stateStore.BasketClassTable().Insert(ctx,
			&api.BasketClass{
				BasketId: basketID,
				ClassId:  class,
			},
		); err != nil {
			return err
		}

		sdkCtx.GasMeter().ConsumeGas(ecocredit.GasCostPerIteration, "ecocredit/basket/MsgCreate class iteration")
	}
	return nil
}
