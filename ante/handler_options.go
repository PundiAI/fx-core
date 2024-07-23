package ante

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	ethante "github.com/evmos/ethermint/app/ante"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper, EVM Keeper and Fee Market Keeper.
type HandlerOptions struct {
	AccountKeeper     AccountKeeper
	BankKeeper        authtypes.BankKeeper
	FeegrantKeeper    FeegrantKeeper
	EvmKeeper         ethante.EVMKeeper
	FeeMarketKeeper   FeeMarketKeeper
	IbcKeeper         *ibckeeper.Keeper
	SignModeHandler   authsigning.SignModeHandler
	SigGasConsumer    ante.SignatureVerificationGasConsumer
	TxFeeChecker      ante.TxFeeChecker
	MaxTxGasWanted    uint64
	InterceptMsgTypes map[int64][]string
	DisabledAuthzMsgs []string
	PendingTxListener ethante.PendingTxListener
}

func (options HandlerOptions) Validate() error {
	if options.AccountKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "account keeper is required for AnteHandler")
	}
	if options.BankKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "bank keeper is required for AnteHandler")
	}
	if options.SignModeHandler == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "sign mode handler is required for ante builder")
	}
	if options.FeeMarketKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "fee market keeper is required for AnteHandler")
	}
	if options.EvmKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "evm keeper is required for AnteHandler")
	}
	return nil
}

func newEthAnteHandler(ctx sdk.Context, options HandlerOptions) sdk.AnteHandler {
	evmParams := options.EvmKeeper.GetParams(ctx)
	evmDenom := evmParams.EvmDenom
	chainID := options.EvmKeeper.ChainID()
	chainCfg := evmParams.GetChainConfig()
	ethCfg := chainCfg.EthereumConfig(chainID)
	baseFee := options.EvmKeeper.GetBaseFee(ctx, ethCfg)
	decorators := []sdk.AnteDecorator{
		ethante.NewEthSetUpContextDecorator(options.EvmKeeper),               // outermost AnteDecorator. SetUpContext must be called first
		ethante.NewEthMempoolFeeDecorator(evmDenom, baseFee),                 // Check eth effective gas price against minimal-gas-prices
		ethante.NewEthMinGasPriceDecorator(options.FeeMarketKeeper, baseFee), // Check eth effective gas price against the global MinGasPrice
		ethante.NewEthValidateBasicDecorator(&evmParams, baseFee),
		ethante.NewEthSigVerificationDecorator(chainID),
		ethante.NewEthAccountVerificationDecorator(options.AccountKeeper, options.EvmKeeper, evmDenom),
		ethante.NewCanTransferDecorator(options.EvmKeeper, baseFee, &evmParams, ethCfg),
		NewEthPubKeyDecorator(options.AccountKeeper),
		ethante.NewEthGasConsumeDecorator(options.EvmKeeper, options.MaxTxGasWanted, ethCfg, evmDenom, baseFee),
		ethante.NewEthIncrementSenderSequenceDecorator(options.AccountKeeper), // innermost AnteDecorator.
		ethante.NewGasWantedDecorator(options.FeeMarketKeeper, ethCfg),
		ethante.NewEthEmitEventDecorator(options.EvmKeeper), // emit eth tx hash and index at the very last ante handler.
	}
	decorators = append(decorators, newTxListenerDecorator(options.PendingTxListener))
	return sdk.ChainAnteDecorators(decorators...)
}

func newCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		// disable the Msg types that cannot be included on an authz.MsgExec msgs field
		ethante.NewAuthzLimiterDecorator(options.DisabledAuthzMsgs),
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		NewPubKeyDecorator(options.AccountKeeper),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		ibcante.NewRedundantRelayDecorator(options.IbcKeeper),
	)
}
