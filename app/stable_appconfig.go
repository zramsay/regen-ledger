//go:build !experimental
// +build !experimental

// DONTCOVER

package app

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	"github.com/cosmos/cosmos-sdk/x/gov"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/regen-network/regen-ledger/types/module/server"
	"github.com/regen-network/regen-ledger/x/ecocredit"
	ecocreditcore "github.com/regen-network/regen-ledger/x/ecocredit/client/core"
	"github.com/regen-network/regen-ledger/x/ecocredit/client/marketplace"
	ecocreditmodule "github.com/regen-network/regen-ledger/x/ecocredit/module"
)

func setCustomModuleBasics() []module.AppModuleBasic {
	return []module.AppModuleBasic{
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distrclient.ProposalHandler,
			upgradeclient.ProposalHandler, upgradeclient.CancelProposalHandler,
			ecocreditcore.CreditTypeProposalHandler, marketplace.AllowDenomProposalHandler,
		),
	}
}

// setCustomModules registers new modules with the server module manager.
// It does nothing here and returns an empty manager since we're not using experimental mode.
func setCustomModules(app *RegenApp, interfaceRegistry types.InterfaceRegistry) *server.Manager {
	return server.NewManager(app.BaseApp, codec.NewProtoCodec(interfaceRegistry))
}
func setCustomKVStoreKeys() []string {
	return []string{}
}

func setCustomMaccPerms() map[string][]string {
	return map[string][]string{}
}

func setCustomOrderBeginBlocker() []string {
	return []string{}
}

func setCustomOrderEndBlocker() []string {
	return []string{}
}

func (app *RegenApp) registerUpgradeHandlers() {

	// mainnet upgrade handler
	const upgradeName = "v4.0.0"
	app.UpgradeKeeper.SetUpgradeHandler(upgradeName, func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		// run state migrations for sdk modules
		toVersion, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
		if err != nil {
			return nil, err
		}

		// run x/ecocredit state migrations
		if err := app.smm.RunMigrations(ctx, app.AppCodec()); err != nil {
			return nil, err
		}
		toVersion[ecocredit.ModuleName] = ecocreditmodule.Module{}.ConsensusVersion()

		// recover funds for community member (regen-1 governance proposal #11)
		if ctx.ChainID() == "regen-1" {
			if err := recoverFunds(ctx, app.AccountKeeper, app.BankKeeper); err != nil {
				return nil, err
			}
		}

		return toVersion, nil
	})

}

func recoverFunds(ctx sdk.Context, ak authkeeper.AccountKeeper, bk bankkeeper.Keeper) error {
	// address with funds inaccessible
	lostAddr, err := sdk.AccAddressFromBech32("regen1c3lpjaq0ytdtsrnjqzmtj3hceavl8fe2vtkj7f")
	if err != nil {
		return err
	}

	// address that the community member has access to
	newAddr, err := sdk.AccAddressFromBech32("regen14tpuqrwf95evu3ejm9z7dn20ttcyzqy3jjpfv4")
	if err != nil {
		return err
	}

	lostAccount := ak.GetAccount(ctx, lostAddr)
	if lostAccount == nil {
		return fmt.Errorf("%s account not found", lostAccount.GetAddress().String())
	}

	newAccount := ak.GetAccount(ctx, newAddr)
	if newAccount == nil {
		return fmt.Errorf("%s account not found", newAccount.GetAddress().String())
	}

	va, ok := lostAccount.(*vestingtypes.PeriodicVestingAccount)
	if !ok {
		return fmt.Errorf("%s is not a vesting account", lostAddr)
	}

	vestingPeriods := va.VestingPeriods
	// unlock vesting tokens
	newVestingPeriods := make([]vestingtypes.Period, len(va.VestingPeriods))
	for i, vp := range va.VestingPeriods {
		vp.Length = 0
		newVestingPeriods[i] = vp
	}
	va.VestingPeriods = newVestingPeriods
	ak.SetAccount(ctx, va)

	// send spendable balance from lost account to new account
	spendable := bk.SpendableCoins(ctx, lostAccount.GetAddress())
	if err := bk.SendCoins(ctx, lostAccount.GetAddress(), newAccount.GetAddress(), spendable); err != nil {
		return err
	}

	newPVA := vestingtypes.NewPeriodicVestingAccount(
		authtypes.NewBaseAccount(newAccount.GetAddress(), newAccount.GetPubKey(), newAccount.GetAccountNumber(), newAccount.GetSequence()),
		va.OriginalVesting, va.StartTime, vestingPeriods,
	)
	ak.SetAccount(ctx, newPVA)

	return nil
}

func (app *RegenApp) setCustomAnteHandler(encCfg simappparams.EncodingConfig, wasmKey *sdk.KVStoreKey, _ *wasm.Config) (sdk.AnteHandler, error) {
	return ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			SignModeHandler: encCfg.TxConfig.SignModeHandler(),
			FeegrantKeeper:  app.FeeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
}

func (app *RegenApp) setCustomModuleManager() []module.AppModule {
	return []module.AppModule{}
}

func (app *RegenApp) setCustomKeepers(_ *baseapp.BaseApp, keys map[string]*sdk.KVStoreKey, appCodec codec.Codec, _ govtypes.Router, _ string,
	_ servertypes.AppOptions,
	_ []wasm.Option) {
}

func setCustomOrderInitGenesis() []string {
	return []string{}
}

func (app *RegenApp) setCustomSimulationManager() []module.AppModuleSimulation {
	return []module.AppModuleSimulation{}
}

func initCustomParamsKeeper(_ *paramskeeper.Keeper) {}

func (app *RegenApp) initializeCustomScopedKeepers() {}
