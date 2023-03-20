// This software is Copyright (c) 2019-2020 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package emoney

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	embank "github.com/MonikaCat/em-ledger/hooks/bank"
	apptypes "github.com/MonikaCat/em-ledger/types"
	"github.com/MonikaCat/em-ledger/x/auth/ante"
	"github.com/MonikaCat/em-ledger/x/authority"
	"github.com/MonikaCat/em-ledger/x/buyback"
	emdistr "github.com/MonikaCat/em-ledger/x/distribution"
	"github.com/MonikaCat/em-ledger/x/inflation"
	"github.com/MonikaCat/em-ledger/x/issuer"
	"github.com/MonikaCat/em-ledger/x/liquidityprovider"
	lptypes "github.com/MonikaCat/em-ledger/x/liquidityprovider/types"
	"github.com/MonikaCat/em-ledger/x/market"
	"github.com/MonikaCat/em-ledger/x/queries"
	queriestypes "github.com/MonikaCat/em-ledger/x/queries/types"
	emslashing "github.com/MonikaCat/em-ledger/x/slashing"
	"github.com/MonikaCat/em-ledger/x/staking"
	historykeeper "github.com/MonikaCat/em-ledger/x/staking/keeper"
	"github.com/MonikaCat/em-ledger/x/upgrade"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/cosmos/ibc-go/v2/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v2/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v2/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v2/modules/core"
	ibcconnectiontypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"
	channelkeeper "github.com/cosmos/ibc-go/v2/modules/core/04-channel/keeper"
	porttypes "github.com/cosmos/ibc-go/v2/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v2/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v2/modules/core/keeper"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/cast"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	db "github.com/tendermint/tm-db"
	dbm "github.com/tendermint/tm-db"
)

const (
	appName = "emoneyd"
)

var (
	DefaultNodeHome = os.ExpandEnv("$HOME/.emd")

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		emdistr.AppModuleBasic{},
		// todo (reviewer) : gov was deactivated in the original app.go but for ICS-20 we should have the `upgradeclient` handlers
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		emslashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		// em modules
		inflation.AppModuleBasic{},
		liquidityprovider.AppModuleBasic{},
		issuer.AppModuleBasic{},
		authority.AppModule{},
		market.AppModule{},
		buyback.AppModule{},
		queries.AppModule{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		emdistr.ModuleName:             nil,
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		// em modules
		inflation.ModuleName:         {authtypes.Minter},
		emslashing.ModuleName:        nil, // TODO Remove this line?
		liquidityprovider.ModuleName: {authtypes.Minter, authtypes.Burner},
		buyback.ModuleName:           {authtypes.Burner},
	}

	// module accounts that are allowed to receive tokens
	allowedReceivingModAcc = map[string]bool{
		emdistr.ModuleName: true,
	}
)

var (
	_ simapp.App              = (*EMoneyApp)(nil)
	_ servertypes.Application = (*EMoneyApp)(nil)
)

type EMoneyApp struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	database     db.DB
	currentBatch db.Batch

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       *embank.ProxyKeeper
	Historykeeper    historykeeper.HistoryKeeper
	CapabilityKeeper *capabilitykeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   emslashing.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IbcKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper
	FeeGrantKeeper   feegrantkeeper.Keeper
	AuthzKeeper      authzkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	// custom modules
	InflationKeeper inflation.Keeper
	LpKeeper        liquidityprovider.Keeper
	IssuerKeeper    issuer.Keeper
	AuthorityKeeper authority.Keeper
	MarketKeeper    *market.Keeper
	BuybackKeeper   buyback.Keeper

	// the module manager
	mm *module.Manager

	// simulation manager
	configurator module.Configurator
}

func (app *EMoneyApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *EMoneyApp) SimulationManager() *module.SimulationManager {
	panic("not supported")
}

type GenesisState map[string]json.RawMessage

// NewApp returns a reference to an initialized App.
func NewApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig EncodingConfig,
	appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *EMoneyApp {
	appCodec := encodingConfig.Marshaler
	legacyAmino := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(appName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		distrtypes.StoreKey, emslashing.StoreKey,
		paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		lptypes.StoreKey, issuer.StoreKey, authority.StoreKey,
		market.StoreKey, market.StoreKeyIdx, buyback.StoreKey,
		inflation.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
	)

	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &EMoneyApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          encodingConfig.TxConfig,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
		database:          createApplicationDatabase(homePath),
	}

	app.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])

	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)

	app.BankKeeper = embank.Wrap(bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.ModuleAccountAddrs(),
	))

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)

	app.AuthzKeeper = authzkeeper.NewKeeper(keys[authzkeeper.StoreKey], appCodec, app.BaseApp.MsgServiceRouter())

	app.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, keys[feegrant.StoreKey], app.AccountKeeper)

	app.Historykeeper = historykeeper.NewHistoryKeeper(appCodec, keys[historykeeper.StoreKey], stakingKeeper, app.database)

	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = emslashing.NewKeeper(
		appCodec, keys[emslashing.StoreKey], &stakingKeeper, app.GetSubspace(emslashing.ModuleName), app.BankKeeper,
		app.database, authtypes.FeeCollectorName,
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath, app.BaseApp)

	{
		/*
		 * This is a test v44 handler
		 */
		const upg44Plan = "v44-upg-test-sample"

		app.UpgradeKeeper.SetUpgradeHandler(
			upg44Plan,
			func(
				ctx sdk.Context, _ upgradetypes.Plan, _ module.VersionMap,
			) (module.VersionMap, error) {
				// set max expected block time parameter. Replace the default with your expected value
				// https://github.com/cosmos/ibc-go/blob/release/v1.0.x/docs/ibc/proto-docs.md#params-2
				// set max block time to 30 seconds
				app.IbcKeeper.ConnectionKeeper.SetParams(
					// should 3-5 minutes in mainnet
					ctx, ibcconnectiontypes.DefaultParams(),
				)

				fromVM := make(map[string]uint64)
				for moduleName := range app.mm.Modules {
					// v40 state is version 1
					// v43 state is version 2
					// v45 state is likely v3
					fromVM[moduleName] = 1
				}

				// EXCEPT Auth needs to run _after_ staking (https://github.com/cosmos/cosmos-sdk/issues/10591),
				// and it seems bank as well (https://github.com/provenance-io/provenance/blob/407c89a7d73854515894161e1526f9623a94c368/app/upgrades.go#L86-L122).
				// So we do this by making auth run last.
				// This is done by setting auth's consensus version to 2, running RunMigrations,
				// then setting it back to 1, and then running migrations again.
				fromVM[authtypes.ModuleName] = 2

				// override versions for _new_ modules as to not skip InitGenesis
				fromVM[authz.ModuleName] = 0
				fromVM[feegrant.ModuleName] = 0

				ctx.Logger().Info("Upgrading to " + upg44Plan)

				newVM, err := app.mm.RunMigrations(ctx, app.configurator, fromVM)
				if err != nil {
					return nil, err
				}

				// now update auth version back to v1, to run auth migration last
				newVM[authtypes.ModuleName] = 1

				return app.mm.RunMigrations(ctx, app.configurator, newVM)
			})

		upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
		if err != nil {
			panic(err)
		}

		if upgradeInfo.Name == upg44Plan && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
			storeUpgrades := store.StoreUpgrades{
				Added: []string{authz.ModuleName, feegrant.ModuleName},
			}

			// configure store loader that checks if version == upgradeHeight and applies store upgrades
			app.SetStoreLoader(
				upgradetypes.UpgradeStoreLoader(
					upgradeInfo.Height, &storeUpgrades,
				),
			)
		}
	}

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.IbcKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName),
		app.Historykeeper, app.UpgradeKeeper, scopedIBCKeeper)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IbcKeeper.ChannelKeeper, &app.IbcKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	app.IbcKeeper.SetRouter(ibcRouter)

	// create evidence keeper with router
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper

	app.InflationKeeper = inflation.NewKeeper(app.appCodec, keys[inflation.StoreKey], app.BankKeeper, app.AccountKeeper, app.StakingKeeper, buyback.AccountName, authtypes.FeeCollectorName)
	app.LpKeeper = liquidityprovider.NewKeeper(app.appCodec, keys[lptypes.StoreKey], app.BankKeeper)
	app.IssuerKeeper = issuer.NewKeeper(app.appCodec, keys[issuer.StoreKey], app.LpKeeper, app.InflationKeeper, app.BankKeeper)
	app.AuthorityKeeper = authority.NewKeeper(app.appCodec, keys[authority.StoreKey], app.IssuerKeeper, app.BankKeeper, app, &app.UpgradeKeeper, app.ParamsKeeper)
	app.MarketKeeper = market.NewKeeper(app.appCodec, keys[market.StoreKey], keys[market.StoreKeyIdx], app.AccountKeeper, app.BankKeeper)
	app.BuybackKeeper = buyback.NewKeeper(app.appCodec, keys[buyback.StoreKey], app.MarketKeeper, app.AccountKeeper, app.StakingKeeper, app.BankKeeper)

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, authsims.RandomGenesisAccounts),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, interfaceRegistry),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		emslashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.Historykeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IbcKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,
		emdistr.NewAppModule(distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper), app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.database),
		liquidityprovider.NewAppModule(app.LpKeeper),
		issuer.NewAppModule(app.IssuerKeeper),
		authority.NewAppModule(app.AuthorityKeeper),
		market.NewAppModule(app.MarketKeeper),
		buyback.NewAppModule(app.BuybackKeeper, app.BankKeeper),
		inflation.NewAppModule(app.InflationKeeper),
		queries.NewAppModule(app.AccountKeeper, app.BankKeeper, app.SlashingKeeper),
	)

	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		//// Cosmos #9800: capability module's begin blocker must come before any modules using capabilities (e.g. IBC)
		capabilitytypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		ibchost.ModuleName,

		authority.ModuleName,
		market.ModuleName,
		inflation.ModuleName,
		emslashing.ModuleName,
		emdistr.ModuleName,
		buyback.ModuleName,

		// Modules that don't seem to use BeginBlock but must now be included for completenes.
		paramstypes.ModuleName,
		banktypes.ModuleName,
		genutiltypes.ModuleName,
		crisistypes.ModuleName,
		issuer.ModuleName,
		vestingtypes.ModuleName,
		feegrant.ModuleName,
		liquidityprovider.ModuleName,
		authtypes.ModuleName,
		authz.ModuleName,
		ibctransfertypes.ModuleName,
		queriestypes.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		stakingtypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,

		// Modules that don't seem to use EndBlock but must now be included for completenes.
		ibchost.ModuleName,
		vestingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		market.ModuleName,
		buyback.ModuleName,
		capabilitytypes.ModuleName,
		upgradetypes.ModuleName,
		evidencetypes.ModuleName,
		liquidityprovider.ModuleName,
		genutiltypes.ModuleName,
		banktypes.ModuleName,
		paramstypes.ModuleName,
		issuer.ModuleName,
		authority.ModuleName,
		inflation.ModuleName,
		queriestypes.ModuleName,
		authtypes.ModuleName,
		slashingtypes.ModuleName,
		distrtypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName, authtypes.ModuleName, banktypes.ModuleName, emdistr.ModuleName, stakingtypes.ModuleName,
		emslashing.ModuleName, crisistypes.ModuleName,
		ibchost.ModuleName, genutiltypes.ModuleName, evidencetypes.ModuleName, ibctransfertypes.ModuleName,
		inflation.ModuleName, issuer.ModuleName, authority.ModuleName, market.ModuleName, buyback.ModuleName,
		liquidityprovider.ModuleName, feegrant.ModuleName, authz.ModuleName,

		// Modules that don't seem to use Genesis initialization but must now be included for completenes.
		upgradetypes.ModuleName, paramstypes.ModuleName, queriestypes.ModuleName, vestingtypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	// TODO hack to initialize application
	// from sdk.bank/module.go: AppModule RegisterServices
	// 	m := keeper.NewMigrator(am.keeper.(keeper.BaseKeeper))
	sdkBk := app.BankKeeper.GetBankKeeper()
	bankMigModule := bank.NewAppModule(appCodec, *sdkBk, app.AccountKeeper)
	proxy := app.mm.Modules[banktypes.ModuleName]
	app.mm.Modules[banktypes.ModuleName] = bankMigModule
	app.mm.RegisterServices(app.configurator)

	app.mm.Modules[banktypes.ModuleName] = proxy

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	anteHandler, err := ante.NewAnteHandler(
		ante.EmAnteHandlerOptions{
			AccountKeeper:    app.AccountKeeper,
			BankKeeper:       app.BankKeeper,
			FeegrantKeeper:   app.FeeGrantKeeper,
			SignModeHandler:  encodingConfig.TxConfig.SignModeHandler(),
			SigGasConsumer:   sdkante.DefaultSigVerificationGasConsumer,
			StakingKeeper:    app.StakingKeeper,
			IBCChannelkeeper: channelkeeper.Keeper{},
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	app.SetAnteHandler(anteHandler)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}

func createApplicationDatabase(rootDir string) db.DB {
	datadirectory := filepath.Join(rootDir, "data")
	emoneydb, err := db.NewGoLevelDB("emoney", datadirectory)
	if err != nil {
		panic(err)
	}

	return emoneydb
}

// load a particular height
func (app *EMoneyApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

func (app *EMoneyApp) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "emz")
}

func (app *EMoneyApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	// Due to how the current Tendermint implementation calculates blocktime, a byzantine 1/3 of voting power can move time forward an arbitrary amount.
	// Have non-malicious nodes shut down if this appears to be happening.
	// This will effectively halt the chain and require off-chain coordination to remedy.
	walltime := time.Now().UTC()
	if walltime.Add(time.Hour).Before(ctx.BlockTime()) {
		s := fmt.Sprintf("Blocktime %v is too far ahead of local wall clock %v.\nSuspending node without processing block %v.\n", ctx.BlockTime(), walltime, ctx.BlockHeight())
		fmt.Println(s)
		panic(s)
	}

	app.currentBatch = app.database.NewBatch() // store in app state as ctx is different in end block
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	ctx = apptypes.WithCurrentBatch(ctx, app.currentBatch)

	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *EMoneyApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	block := ctx.BlockHeader()
	proposerAddress := block.GetProposerAddress()
	app.Logger(ctx).Info(fmt.Sprintf("Endblock: Block %v was proposed by %v", ctx.BlockHeight(), sdk.ValAddress(proposerAddress)))

	response := app.mm.EndBlock(ctx, req)
	err := app.currentBatch.Write() // Write non-IAVL state to database
	if err != nil {                 // todo (reviewer): should we panic or ignore? panics are not handled downstream will cause a crash
		panic(err)
	}
	return response
}

// InitChainer application update at chain initialization
func (app *EMoneyApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())

	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

func (app *EMoneyApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// AppCodec returns SimApp's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *EMoneyApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns SimApp's InterfaceRegistry
func (app *EMoneyApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *EMoneyApp) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *EMoneyApp) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *EMoneyApp) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *EMoneyApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *EMoneyApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(clientCtx, apiSvr.Router)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *EMoneyApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *EMoneyApp) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// RegisterSwaggerAPI registers swagger route with API Server
func RegisterSwaggerAPI(ctx client.Context, rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

// GetMaccs returns a copy of the module accounts
func GetMaccs() map[string]bool {
	maccs := make(map[string]bool)
	for k := range maccPerms {
		maccs[k] = true
	}
	return maccs
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(emslashing.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)

	return paramsKeeper
}

func (app EMoneyApp) SetMinimumGasPrices(gasPricesStr string) (err error) {
	if _, err = sdk.ParseDecCoins(gasPricesStr); err != nil {
		return
	}

	baseapp.SetMinGasPrices(gasPricesStr)(app.BaseApp)
	return
}

func init() {
	sdk.DefaultPowerReduction = sdk.OneInt()
}
