package keys_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/functionx/fx-core/v7/app"
	"github.com/functionx/fx-core/v7/client/cli/keys"
)

func TestNewKeyOutput(t *testing.T) {
	pubKeyJson := `{"@type":"/ethermint.crypto.v1.ethsecp256k1.PubKey","key":"A4MYU1tUEF1Keq5gwI/EX5aHGBtP38YlvRp1P6c5f+11"}`
	encodingConfig := app.MakeEncodingConfig()
	var pubKey cryptotypes.PubKey
	err := encodingConfig.Codec.UnmarshalInterfaceJSON([]byte(pubKeyJson), &pubKey)
	assert.NoError(t, err)
	address := sdk.AccAddress(pubKey.Address().Bytes())
	keyOutput, err := keys.NewKeyOutput("test", keyring.TypeLocal, address, pubKey)
	assert.NoError(t, err)
	assert.Equal(t, "cosmos17w0adeg64ky0daxwd2ugyuneellmjgnxramjtq", keyOutput.Address)
	assert.Equal(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", keyOutput.Eip55Address)
	assert.Equal(t, "", keyOutput.Mnemonic)
	assert.Equal(t, "test", keyOutput.Name)
	assert.Equal(t, pubKeyJson, keyOutput.PubKey)
	assert.Equal(t, "local", keyOutput.Type)
}
