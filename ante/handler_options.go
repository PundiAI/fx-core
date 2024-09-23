package ante

import (
	errorsmod "cosmossdk.io/errors"
	txsigning "cosmossdk.io/x/tx/signing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethante "github.com/evmos/ethermint/app/ante"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

// HandlerOptions extend the SDK's AnteHandler options by requiring the IBC
// channel keeper, EVM Keeper and Fee Market Keeper.
type HandlerOptions struct {
	AccountKeeper     evmtypes.AccountKeeper
	BankKeeper        authtypes.BankKeeper
	FeegrantKeeper    ante.FeegrantKeeper
	EvmKeeper         ethante.EVMKeeper
	FeeMarketKeeper   evmtypes.FeeMarketKeeper
	IbcKeeper         *ibckeeper.Keeper
	GovKeeper         Govkeeper
	SignModeHandler   *txsigning.HandlerMap
	SigGasConsumer    ante.SignatureVerificationGasConsumer
	TxFeeChecker      ante.TxFeeChecker
	MaxTxGasWanted    uint64
	InterceptMsgTypes map[int64][]string
	DisabledAuthzMsgs []string
	PendingTxListener ethante.PendingTxListener

	UnsafeUnorderedTx bool
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
	if options.IbcKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "ibc keeper is required for AnteHandler")
	}
	if options.GovKeeper == nil {
		return errorsmod.Wrap(errortypes.ErrLogic, "gov keeper is required for AnteHandler")
	}
	return nil
}

func newEthAnteHandler(options HandlerOptions) sdk.AnteHandler {
	decorators := []sdk.AnteDecorator{
		NewEthPubKeyDecorator(options.AccountKeeper),
		newTxListenerDecorator(options.PendingTxListener),
	}
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		blockCfg, err := options.EvmKeeper.EVMBlockConfig(ctx, options.EvmKeeper.ChainID())
		if err != nil {
			return ctx, errorsmod.Wrap(errortypes.ErrLogic, err.Error())
		}
		evmParams := &blockCfg.Params
		evmDenom := evmParams.EvmDenom
		feemarketParams := &blockCfg.FeeMarketParams
		baseFee := blockCfg.BaseFee
		rules := blockCfg.Rules

		// all transactions must implement FeeTx
		_, ok := tx.(sdk.FeeTx)
		if !ok {
			return ctx, errorsmod.Wrapf(errortypes.ErrInvalidType, "invalid transaction type %T, expected sdk.FeeTx", tx)
		}

		// We need to setup an empty gas config so that the gas is consistent with Ethereum.
		ctx, err = ethante.SetupEthContext(ctx)
		if err != nil {
			return ctx, err
		}

		if err = ethante.CheckEthMempoolFee(ctx, tx, simulate, baseFee, evmDenom); err != nil {
			return ctx, err
		}

		if err = ethante.CheckEthMinGasPrice(tx, feemarketParams.MinGasPrice, baseFee); err != nil {
			return ctx, err
		}

		if err = ethante.ValidateEthBasic(ctx, tx, evmParams, baseFee); err != nil {
			return ctx, err
		}

		ethSigner := ethtypes.MakeSigner(blockCfg.ChainConfig, blockCfg.BlockNumber)
		err = ethante.VerifyEthSig(tx, ethSigner)
		if err != nil {
			return ctx, err
		}

		// AccountGetter cache the account objects during the ante handler execution,
		// it's safe because there's no store branching in the ante handlers.
		accountGetter := ethante.NewCachedAccountGetter(ctx, options.AccountKeeper)

		if err = ethante.VerifyEthAccount(ctx, tx, options.EvmKeeper, evmDenom, accountGetter); err != nil {
			return ctx, err
		}

		if err = ethante.CheckEthCanTransfer(ctx, tx, baseFee, rules, options.EvmKeeper, evmParams); err != nil {
			return ctx, err
		}

		ctx, err = ethante.CheckEthGasConsume(
			ctx, tx, rules, options.EvmKeeper,
			baseFee, options.MaxTxGasWanted, evmDenom,
		)
		if err != nil {
			return ctx, err
		}

		if err = ethante.CheckAndSetEthSenderNonce(ctx, tx, options.AccountKeeper, options.UnsafeUnorderedTx, accountGetter); err != nil {
			return ctx, err
		}

		if len(decorators) > 0 {
			return sdk.ChainAnteDecorators(decorators...)(ctx, tx, simulate)
		}
		return ctx, nil
	}
}

func newCosmosAnteHandler(options HandlerOptions) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ethante.RejectMessagesDecorator{}, // reject MsgEthereumTxs
		// disable the Msg types that cannot be included on an authz.MsgExec msgs field
		NewDisableMsgDecorator(options.DisabledAuthzMsgs, options.GovKeeper),
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
