package cmd_test

import (
	"fmt"
	"testing"

	emoney "github.com/MonikaCat/em-ledger"

	"github.com/stretchr/testify/require"

	"github.com/MonikaCat/em-ledger/cmd/emd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
)

func TestInitCmd(t *testing.T) {
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",       // Test the init cmd
		"emapp-test", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
	})

	require.NoError(t, svrcmd.Execute(rootCmd, emoney.DefaultNodeHome))
}
