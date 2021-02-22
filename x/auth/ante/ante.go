package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak AccountKeeper, bankKeeper types.BankKeeper, feeGrantKeeper feegrantkeeper.Keeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		NewMempoolFeeDecorator(),
		NewValidateBasicDecorator(),
		TxTimeoutHeightDecorator{},
		NewValidateMemoDecorator(ak),
		NewConsumeGasForTxSizeDecorator(ak),
		NewDeductGrantedFeeDecorator(ak, bankKeeper, feeGrantKeeper),
		// NewRejectFeeGranterDecorator(),
		NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		NewValidateSigCountDecorator(ak),
		// NewDeductFeeDecorator(ak, bankKeeper),
		NewSigGasConsumeDecorator(ak, sigGasConsumer),
		NewSigVerificationDecorator(ak, signModeHandler),
		NewIncrementSequenceDecorator(ak),
	)
}
