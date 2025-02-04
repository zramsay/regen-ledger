package marketplace

import (
	"strconv"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/regen-network/gocuke"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/orm/types/ormerrors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	api "github.com/regen-network/regen-ledger/api/regen/ecocredit/marketplace/v1"
	coreapi "github.com/regen-network/regen-ledger/api/regen/ecocredit/v1"
	"github.com/regen-network/regen-ledger/x/ecocredit/marketplace"
)

type cancelSellOrder struct {
	*baseSuite
	alice            sdk.AccAddress
	bob              sdk.AccAddress
	creditTypeAbbrev string
	classId          string
	batchDenom       string
	sellOrderId      uint64
	askPrice         *sdk.Coin
	quantity         string
	res              *marketplace.MsgCancelSellOrderResponse
	err              error
}

func TestCancelSellOrder(t *testing.T) {
	gocuke.NewRunner(t, &cancelSellOrder{}).Path("./features/msg_cancel_sell_order.feature").Run()
}

func (s *cancelSellOrder) Before(t gocuke.TestingT) {
	s.baseSuite = setupBase(t)
	s.alice = s.addr
	s.bob = s.addr2
	s.creditTypeAbbrev = "C"
	s.classId = "C01"
	s.batchDenom = "C01-001-20200101-20210101-001"
	s.askPrice = &sdk.Coin{
		Denom:  "regen",
		Amount: sdk.NewInt(100),
	}
	s.quantity = "100"
}

func (s *cancelSellOrder) AliceCreatedASellOrderWithId(a string) {
	id, err := strconv.ParseUint(a, 10, 32)
	require.NoError(s.t, err)

	s.sellOrderId = id

	s.sellOrderSetup()
}

func (s *cancelSellOrder) AliceCreatedASellOrderWithIdAndQuantity(a string, b string) {
	id, err := strconv.ParseUint(a, 10, 32)
	require.NoError(s.t, err)

	s.sellOrderId = id
	s.quantity = b

	s.sellOrderSetup()
}

func (s *cancelSellOrder) AliceHasTheBatchBalance(a gocuke.DocString) {
	balance := &coreapi.BatchBalance{}
	err := jsonpb.UnmarshalString(a.Content, balance)
	require.NoError(s.t, err)

	batch, err := s.coreStore.BatchTable().GetByDenom(s.ctx, s.batchDenom)
	require.NoError(s.t, err)

	balance.BatchKey = batch.Key
	balance.Address = s.alice

	// Save because the balance already exists from sellOrderSetup
	err = s.coreStore.BatchBalanceTable().Save(s.ctx, balance)
	require.NoError(s.t, err)
}

func (s *cancelSellOrder) AliceAttemptsToCancelTheSellOrderWithId(a string) {
	id, err := strconv.ParseUint(a, 10, 32)
	require.NoError(s.t, err)

	s.res, s.err = s.k.CancelSellOrder(s.ctx, &marketplace.MsgCancelSellOrder{
		Seller:      s.alice.String(),
		SellOrderId: id,
	})
}

func (s *cancelSellOrder) BobAttemptsToCancelTheSellOrderWithId(a string) {
	id, err := strconv.ParseUint(a, 10, 32)
	require.NoError(s.t, err)

	s.res, s.err = s.k.CancelSellOrder(s.ctx, &marketplace.MsgCancelSellOrder{
		Seller:      s.bob.String(),
		SellOrderId: id,
	})
}

func (s *cancelSellOrder) ExpectNoError() {
	require.NoError(s.t, s.err)
}

func (s *cancelSellOrder) ExpectTheError(a string) {
	require.EqualError(s.t, s.err, a)
}

func (s *cancelSellOrder) ExpectAliceBatchBalance(a gocuke.DocString) {
	expected := &coreapi.BatchBalance{}
	err := jsonpb.UnmarshalString(a.Content, expected)
	require.NoError(s.t, err)

	batch, err := s.coreStore.BatchTable().GetByDenom(s.ctx, s.batchDenom)
	require.NoError(s.t, err)

	balance, err := s.coreStore.BatchBalanceTable().Get(s.ctx, s.alice, batch.Key)
	require.NoError(s.t, err)

	require.Equal(s.t, expected.RetiredAmount, balance.RetiredAmount)
	require.Equal(s.t, expected.TradableAmount, balance.TradableAmount)
	require.Equal(s.t, expected.EscrowedAmount, balance.EscrowedAmount)
}

func (s *cancelSellOrder) ExpectNoSellOrderWithId(a string) {
	id, err := strconv.ParseUint(a, 10, 32)
	require.NoError(s.t, err)

	_, err = s.marketStore.SellOrderTable().Get(s.ctx, id)
	require.ErrorContains(s.t, err, ormerrors.NotFound.Error())
}

func (s *cancelSellOrder) sellOrderSetup() {
	err := s.coreStore.ClassTable().Insert(s.ctx, &coreapi.Class{
		Id:               s.classId,
		CreditTypeAbbrev: s.creditTypeAbbrev,
	})
	require.NoError(s.t, err)

	batchKey, err := s.coreStore.BatchTable().InsertReturningID(s.ctx, &coreapi.Batch{
		Denom: s.batchDenom,
	})
	require.NoError(s.t, err)

	err = s.coreStore.BatchBalanceTable().Insert(s.ctx, &coreapi.BatchBalance{
		BatchKey:       batchKey,
		Address:        s.alice,
		EscrowedAmount: s.quantity,
	})
	require.NoError(s.t, err)

	err = s.coreStore.BatchSupplyTable().Insert(s.ctx, &coreapi.BatchSupply{
		BatchKey:       batchKey,
		TradableAmount: s.quantity,
	})
	require.NoError(s.t, err)

	marketKey, err := s.marketStore.MarketTable().InsertReturningID(s.ctx, &api.Market{
		CreditTypeAbbrev: s.creditTypeAbbrev,
		BankDenom:        s.askPrice.Denom,
	})
	require.NoError(s.t, err)

	sellOrderId, err := s.marketStore.SellOrderTable().InsertReturningID(s.ctx, &api.SellOrder{
		Seller:    s.alice,
		BatchKey:  batchKey,
		Quantity:  s.quantity,
		MarketId:  marketKey,
		AskAmount: s.askPrice.Amount.String(),
	})
	require.NoError(s.t, err)
	require.Equal(s.t, sellOrderId, s.sellOrderId)
}
