package cmd

import (
	"github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/spf13/cobra"

	fxkeys "github.com/pundiai/fx-core/v8/client/cli/keys"
)

// Commands registers a sub-tree of commands to interact with
// local private key storage.
func keyCommands(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage your application's keys",
		Long: `Keyring management commands. These keys may be in any format supported by the
Tendermint crypto library and can be used by light-clients, full nodes, or any other application
that needs to sign with a private key.

The keyring supports the following backends:

    os          Uses the operating system's default credentials store.
    file        Uses encrypted file-based keystore within the app's configuration directory.
                This keyring will request a password each time it is accessed, which may occur
                multiple times in a single command resulting in repeated password prompts.
    kwallet     Uses KDE Wallet Manager as a credentials management application.
    pass        Uses the pass command line utility to store and retrieve keys.
    test        Stores keys insecurely to disk. It does not prompt for a password to be unlocked
                and it should be use only for testing purposes.

kwallet and pass backends depend on external tools. Refer to their respective documentation for more
information:
    KWallet     https://github.com/KDE/kwallet
    pass        https://www.passwordstore.org/

The pass backend requires GnuPG: https://gnupg.org/
`,
	}

	cmd.AddCommand(
		keys.MnemonicKeyCommand(),
		fxkeys.AddKeyCommand(),
		fxkeys.ExportKeyCommand(),
		fxkeys.ImportKeyCommand(),
		fxkeys.ListKeysCmd(),
		fxkeys.ShowKeysCmd(),
		fxkeys.ParseAddressCommand(),
		keys.ListKeyTypesCmd(),
		keys.DeleteKeyCommand(),
		keys.RenameKeyCommand(),
		keys.ParseKeyStringCommand(),
		keys.MigrateCommand(),
	)

	cmd.PersistentFlags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	cmd.PersistentFlags().String(flags.FlagKeyringDir, "", "The client Keyring directory; if omitted, the default 'home' directory will be used")
	cmd.PersistentFlags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.PersistentFlags().StringP(cli.OutputFlag, "o", "text", "Output format (text|json)")

	return cmd
}
