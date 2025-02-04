package testsuite

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/regen-network/regen-ledger/types"
	"github.com/regen-network/regen-ledger/types/testutil/cli"
	coreclient "github.com/regen-network/regen-ledger/x/ecocredit/client"
	marketplaceclient "github.com/regen-network/regen-ledger/x/ecocredit/client/marketplace"
	"github.com/regen-network/regen-ledger/x/ecocredit/core"
	"github.com/regen-network/regen-ledger/x/ecocredit/marketplace"
)

func (s *IntegrationTestSuite) TestQueryClassesCmd() {
	val := s.network.Validators[0]
	val2 := s.network.Validators[1]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	classId, err := s.createClass(clientCtx, &core.MsgCreateClass{
		Admin:            val.Address.String(),
		Issuers:          []string{val.Address.String()},
		Metadata:         "metadata",
		CreditTypeAbbrev: validCreditTypeAbbrev,
		Fee:              &core.DefaultParams().CreditClassFee[0],
	})
	s.Require().NoError(err)
	classId2, err := s.createClass(clientCtx, &core.MsgCreateClass{
		Admin:            val.Address.String(),
		Issuers:          []string{val.Address.String(), val2.Address.String()},
		Metadata:         "metadata2",
		CreditTypeAbbrev: validCreditTypeAbbrev,
		Fee:              &core.DefaultParams().CreditClassFee[0],
	})
	s.Require().NoError(err)
	classIds := [2]string{classId, classId2}

	testCases := []struct {
		name               string
		args               []string
		expectErr          bool
		expectedErrMsg     string
		expectedAmtClasses int
	}{
		{
			name:           "too many args",
			args:           []string{"abcde"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 0 arg(s), received 1",
		},
		{
			name:               "no pagination flags",
			args:               []string{},
			expectErr:          false,
			expectedAmtClasses: -1,
		},
		{
			name: "limit 1",
			args: []string{
				fmt.Sprintf("--%s=1", flags.FlagLimit),
			},
			expectErr:          false,
			expectedAmtClasses: 1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryClassesCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryClassesResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))

				if tc.expectedAmtClasses > 0 {
					s.Require().Len(res.Classes, tc.expectedAmtClasses)
				} else {
					resClassIds := make([]string, len(res.Classes))
					for i, cls := range res.Classes {
						resClassIds[i] = cls.Id
					}
					for _, id := range classIds {
						s.Require().Contains(resClassIds, id)
					}
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryClassCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	class := &core.MsgCreateClass{
		Admin:            val.Address.String(),
		Issuers:          []string{val.Address.String()},
		Metadata:         "hi",
		CreditTypeAbbrev: validCreditTypeAbbrev,
		Fee:              &core.DefaultParams().CreditClassFee[0],
	}
	classId, err := s.createClass(clientCtx, class)
	s.Require().NoError(err)

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
		expectedClass  *core.ClassInfo
	}{
		{
			name:           "missing args",
			args:           []string{},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:           "too many args",
			args:           []string{"abcde", "abcde"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name:      "valid credit class",
			args:      []string{classId},
			expectErr: false,
			expectedClass: &core.ClassInfo{
				Id:               classId,
				Admin:            val.Address.String(),
				Metadata:         class.Metadata,
				CreditTypeAbbrev: class.CreditTypeAbbrev,
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryClassCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryClassResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().Equal(tc.expectedClass, res.Class)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryBatchesCmd() {
	ctx := s.val.ClientCtx
	ctx.OutputFormat = "JSON"

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:           "too many args",
			args:           []string{"foo"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 0 arg(s), received 1",
		},
		{
			name: "valid with pagination",
			args: []string{
				fmt.Sprintf("--%s", flags.FlagCountTotal),
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryBatchesCmd()
			out, err := cli.ExecTestCLICmd(ctx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryBatchesResponse
				s.Require().NoError(ctx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().True(len(res.Batches) > 0)
				s.Require().NotNil(res.Pagination)
				s.Require().True(res.Pagination.Total > 0)
				denoms := make([]string, len(res.Batches))
				for i, batch := range res.Batches {
					denoms[i] = batch.Denom
				}
				s.Require().Contains(denoms, s.batchDenom)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryBatchesByIssuerCmd() {
	ctx := s.val.ClientCtx
	ctx.OutputFormat = "JSON"

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:           "missing args",
			args:           []string{},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:           "too many args",
			args:           []string{"foo", "bar"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name: "valid with pagination",
			args: []string{
				s.addr1.String(),
				fmt.Sprintf("--%s", flags.FlagCountTotal),
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryBatchesByIssuerCmd()
			out, err := cli.ExecTestCLICmd(ctx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryBatchesByIssuerResponse
				s.Require().NoError(ctx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().True(len(res.Batches) > 0)
				s.Require().NotNil(res.Pagination)
				s.Require().True(res.Pagination.Total > 0)
				denoms := make([]string, len(res.Batches))
				for i, batch := range res.Batches {
					denoms[i] = batch.Denom
				}
				s.Require().Contains(denoms, s.batchDenom)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryBatchesByClassCmd() {
	ctx := s.val.ClientCtx
	ctx.OutputFormat = "JSON"

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:           "missing args",
			args:           []string{},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:           "too many args",
			args:           []string{"foo", "bar"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name: "valid with pagination",
			args: []string{
				s.classId,
				fmt.Sprintf("--%s", flags.FlagCountTotal),
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryBatchesByClassCmd()
			out, err := cli.ExecTestCLICmd(ctx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryBatchesByClassResponse
				s.Require().NoError(ctx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().True(len(res.Batches) > 0)
				s.Require().NotNil(res.Pagination)
				s.Require().True(res.Pagination.Total > 0)
				denoms := make([]string, len(res.Batches))
				for i, batch := range res.Batches {
					denoms[i] = batch.Denom
				}
				s.Require().Contains(denoms, s.batchDenom)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryBatchesByProjectCmd() {
	ctx := s.val.ClientCtx
	ctx.OutputFormat = "JSON"

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:           "missing args",
			args:           []string{},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:           "too many args",
			args:           []string{"foo", "bar"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name: "valid with pagination",
			args: []string{
				s.projectId,
				fmt.Sprintf("--%s", flags.FlagCountTotal),
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryBatchesByProjectCmd()
			out, err := cli.ExecTestCLICmd(ctx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryBatchesByProjectResponse
				s.Require().NoError(ctx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().True(len(res.Batches) > 0)
				s.Require().NotNil(res.Pagination)
				s.Require().True(res.Pagination.Total > 0)
				denoms := make([]string, len(res.Batches))
				for i, batch := range res.Batches {
					denoms[i] = batch.Denom
				}
				s.Require().Contains(denoms, s.batchDenom)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryBatchCmd() {
	ctx := s.val.ClientCtx
	ctx.OutputFormat = "JSON"

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:           "missing args",
			args:           []string{},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:           "too many args",
			args:           []string{"abcde", "abcde"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name:      "valid credit batch",
			args:      []string{s.batchDenom},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryBatchCmd()
			out, err := cli.ExecTestCLICmd(ctx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryBatchResponse
				s.Require().NoError(ctx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().Equal(res.Batch.Denom, s.batchDenom)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryBalanceCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	_, _, batchDenom := s.createClassProjectBatch(clientCtx, val.Address.String())

	testCases := []struct {
		name                   string
		args                   []string
		expectErr              bool
		expectedErrMsg         string
		expectedTradableAmount string
		expectedRetiredAmount  string
	}{
		{
			name:           "missing args",
			args:           []string{},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 2 arg(s), received 0",
		},
		{
			name:           "too many args",
			args:           []string{"abcde", "abcde", "abcde"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 2 arg(s), received 3",
		},
		{
			name:                   "valid",
			args:                   []string{batchDenom, val.Address.String()},
			expectErr:              false,
			expectedTradableAmount: "100",
			expectedRetiredAmount:  "0.000001",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryBalanceCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryBalanceResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().Equal(res.Balance.Address, val.Address.String())
				s.Require().NotEmpty(res.Balance.TradableAmount)
				s.Require().NotEmpty(res.Balance.RetiredAmount)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQuerySupplyCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	_, _, batchDenom := s.createClassProjectBatch(clientCtx, val.Address.String())

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:           "missing args",
			args:           []string{},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:           "too many args",
			args:           []string{"abcde", "abcde"},
			expectErr:      true,
			expectedErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name:      "valid credit batch",
			args:      []string{batchDenom},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QuerySupplyCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QuerySupplyResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().NotEmpty(res.TradableSupply)
				s.Require().NotEmpty(res.RetiredSupply)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryCreditTypesCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
	}{
		{
			name:           "valid",
			args:           []string{},
			expectErr:      false,
			expectedErrMsg: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryCreditTypesCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expectedErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryCreditTypesResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().Greater(len(res.CreditTypes), 0)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryParamsCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	require := s.Require()

	cmd := coreclient.QueryParamsCmd()
	out, err := cli.ExecTestCLICmd(clientCtx, cmd, []string{})
	require.NoError(err)

	var params core.QueryParamsResponse
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &params))
	require.NoError(err)

	require.Equal(core.DefaultParams(), *params.Params)
}

func (s *IntegrationTestSuite) TestQuerySellOrderCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	_, _, batchDenom := s.createClassProjectBatch(clientCtx, val.Address.String())
	validAsk := sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)
	expiration, err := types.ParseDate("expiration", "2050-03-11")
	s.Require().NoError(err)
	orderIds, err := s.createSellOrder(clientCtx, &marketplace.MsgSell{
		Seller: val.Address.String(),
		Orders: []*marketplace.MsgSell_Order{
			{BatchDenom: batchDenom, Quantity: "10", AskPrice: &validAsk, Expiration: &expiration},
		},
	})
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "missing args",
			args:      []string{},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:      "too many args",
			args:      []string{"foo", "bar"},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name:      "valid",
			args:      []string{fmt.Sprintf("%d", orderIds[0])},
			expErr:    false,
			expErrMsg: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := marketplaceclient.QuerySellOrderCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res marketplace.QuerySellOrderResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().Equal(res.SellOrder.Seller, val.Address.String())
				s.Require().Equal(res.SellOrder.Quantity, "10")
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQuerySellOrdersCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	_, _, batchDenom := s.createClassProjectBatch(clientCtx, val.Address.String())
	validAsk := sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)
	expiration, err := types.ParseDate("expiration", "2050-03-11")
	s.Require().NoError(err)

	_, err = s.createSellOrder(clientCtx, &marketplace.MsgSell{
		Seller: val.Address.String(),
		Orders: []*marketplace.MsgSell_Order{
			{BatchDenom: batchDenom, Quantity: "10", AskPrice: &validAsk, Expiration: &expiration},
			{BatchDenom: batchDenom, Quantity: "3", AskPrice: &validAsk, Expiration: &expiration},
		},
	})
	s.Require().NoError(err)
	testCases := []struct {
		name      string
		args      []string
		expErr    bool
		expErrMsg string
		expOrders []*marketplace.SellOrder
	}{
		{
			name:      "too many args",
			args:      []string{"foo"},
			expErr:    true,
			expErrMsg: "Error: accepts 0 arg(s), received 1",
		},
		{
			name:      "valid",
			args:      []string{fmt.Sprintf("--%s", flags.FlagCountTotal)},
			expErr:    false,
			expErrMsg: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := marketplaceclient.QuerySellOrdersCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res marketplace.QuerySellOrdersResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().NotNil(res.Pagination)
				s.Require().True(res.Pagination.Total > 1)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQuerySellOrdersBySellerCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	_, _, batchDenom := s.createClassProjectBatch(clientCtx, val.Address.String())
	validAsk := sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)
	expiration, err := types.ParseDate("expiration", "2050-03-11")
	s.Require().NoError(err)
	_, err = s.createSellOrder(clientCtx, &marketplace.MsgSell{
		Seller: val.Address.String(),
		Orders: []*marketplace.MsgSell_Order{
			{BatchDenom: batchDenom, Quantity: "10", AskPrice: &validAsk, Expiration: &expiration},
			{BatchDenom: batchDenom, Quantity: "3", AskPrice: &validAsk, Expiration: &expiration},
		},
	})
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "missing args",
			args:      []string{},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:      "too many args",
			args:      []string{"foo", "bar"},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name:      "valid",
			args:      []string{val.Address.String(), fmt.Sprintf("--%s", flags.FlagCountTotal)},
			expErr:    false,
			expErrMsg: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := marketplaceclient.QuerySellOrdersBySellerCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res marketplace.QuerySellOrdersBySellerResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().NotNil(res.Pagination)
				s.Require().True(res.Pagination.Total > 1)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQuerySellOrdersByBatchDenomCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	_, _, batchDenom := s.createClassProjectBatch(clientCtx, val.Address.String())
	validAsk := sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)
	expiration, err := types.ParseDate("expiration", "2050-03-11")
	s.Require().NoError(err)

	_, err = s.createSellOrder(clientCtx, &marketplace.MsgSell{
		Seller: val.Address.String(),
		Orders: []*marketplace.MsgSell_Order{
			{BatchDenom: batchDenom, Quantity: "10", AskPrice: &validAsk, Expiration: &expiration},
			{BatchDenom: batchDenom, Quantity: "3", AskPrice: &validAsk, Expiration: &expiration},
		},
	})
	s.Require().NoError(err)

	testCases := []struct {
		name      string
		args      []string
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "missing args",
			args:      []string{},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:      "too many args",
			args:      []string{"foo", "bar"},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name:      "valid",
			args:      []string{batchDenom, fmt.Sprintf("--%s", flags.FlagCountTotal)},
			expErr:    false,
			expErrMsg: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := marketplaceclient.QuerySellOrdersByBatchDenomCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res marketplace.QuerySellOrdersByBatchDenomResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().NotNil(res.Pagination)
				s.Require().Len(res.SellOrders, 2)
				s.Require().Equal(uint64(2), res.Pagination.Total)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryProjectsCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	classId, err := s.createClass(clientCtx, &core.MsgCreateClass{
		Admin:            val.Address.String(),
		Issuers:          []string{val.Address.String()},
		Metadata:         "foo",
		CreditTypeAbbrev: validCreditTypeAbbrev,
		Fee:              &core.DefaultParams().CreditClassFee[0],
	})
	s.Require().NoError(err)
	pID, err := s.createProject(clientCtx, &core.MsgCreateProject{
		Issuer:       val.Address.String(),
		ClassId:      classId,
		Metadata:     "foo",
		Jurisdiction: "US-OR",
	})
	s.Require().NoError(err)
	pID2, err := s.createProject(clientCtx, &core.MsgCreateProject{
		Issuer:       val.Address.String(),
		ClassId:      classId,
		Metadata:     "foo",
		Jurisdiction: "US-OR",
	})
	s.Require().NoError(err)
	projectIds := [2]string{pID, pID2}
	testCases := []struct {
		name      string
		args      []string
		expErr    bool
		expErrMsg string
		expLen    int
	}{
		{
			name:      "no args",
			args:      []string{},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:      "too many args",
			args:      []string{"foo", "bar"},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name:   "valid query",
			args:   []string{classId, fmt.Sprintf("--%s", flags.FlagCountTotal)},
			expErr: false,
			expLen: 2,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryProjectsCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				s.Require().Error(err)
				s.Require().Contains(out.String(), tc.expErrMsg)
			} else {
				s.Require().NoError(err, out.String())

				var res core.QueryProjectsResponse
				s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				s.Require().Len(res.Projects, tc.expLen)
				s.Require().Equal(res.Pagination.Total, uint64(2))
				for _, project := range res.Projects {
					s.Require().Contains(projectIds, project.Id)
				}
			}
		})
	}

}

func (s *IntegrationTestSuite) TestQueryProjectCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	require := s.Require()
	_, projectId, _ := s.createClassProjectBatch(clientCtx, val.Address.String())

	testCases := []struct {
		name      string
		args      []string
		expErr    bool
		expErrMsg string
	}{
		{
			name:      "no args",
			args:      []string{},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 0",
		},
		{
			name:      "too many args",
			args:      []string{"foo", "bar"},
			expErr:    true,
			expErrMsg: "Error: accepts 1 arg(s), received 2",
		},
		{
			name:      "valid query",
			args:      []string{projectId},
			expErr:    false,
			expErrMsg: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryProjectCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				require.Error(err)
				require.Contains(out.String(), tc.expErrMsg)
			} else {
				require.NoError(err, out.String())

				var res core.QueryProjectResponse
				require.NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				require.Equal(projectId, res.Project.Id)
			}
		})
	}

}

func (s *IntegrationTestSuite) TestQueryClassIssuersCmd() {
	val := s.network.Validators[0]
	val2 := s.network.Validators[1]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"
	require := s.Require()

	classId, err := s.createClass(clientCtx, &core.MsgCreateClass{
		Admin:            val.Address.String(),
		Issuers:          []string{val.Address.String(), val2.Address.String()},
		Metadata:         "metadata",
		CreditTypeAbbrev: validCreditTypeAbbrev,
		Fee:              &core.DefaultParams().CreditClassFee[0],
	})
	require.NoError(err)

	testCases := []struct {
		name           string
		args           []string
		expectErr      bool
		expectedErrMsg string
		numItems       int
	}{
		{
			name:           "no pagination flags",
			args:           []string{classId},
			expectErr:      false,
			expectedErrMsg: "",
			numItems:       -1,
		},
		{
			name:           "pagination limit 1",
			args:           []string{classId, "--limit=1"},
			expectErr:      false,
			expectedErrMsg: "",
			numItems:       1,
		},
		{
			name:           "class not found",
			args:           []string{"Z100"},
			expectErr:      true,
			expectedErrMsg: "not found",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := coreclient.QueryClassIssuersCmd()
			out, err := cli.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				require.Error(err)
				require.Contains(out.String(), tc.expectedErrMsg)
			} else {
				require.NoError(err, out.String())

				var res core.QueryClassIssuersResponse
				require.NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
				if tc.numItems > 0 {
					require.Len(res.Issuers, tc.numItems)
				} else {
					require.GreaterOrEqual(len(res.Issuers), 1)
				}
			}
		})
	}
}

func (s *IntegrationTestSuite) TestQueryAllowedDenomsCmd() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	clientCtx.OutputFormat = "JSON"

	cmd := marketplaceclient.QueryAllowedDenomsCmd()
	out, err := cli.ExecTestCLICmd(clientCtx, cmd, []string{fmt.Sprintf("--%s", flags.FlagCountTotal)})
	s.Require().NoError(err)

	var res marketplace.QueryAllowedDenomsResponse
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(out.Bytes(), &res))
	for _, d := range res.AllowedDenoms {
		s.Require().Contains(s.allowedDenoms, d.BankDenom)
	}
	s.Require().Equal(uint64(len(s.allowedDenoms)), res.Pagination.Total)
}
