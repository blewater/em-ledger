package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/tendermint/libs/cli"
)

const (
	flagClientHome = "home-client"
)

// Adapted from Cosmos-sdk : x/genaccounts/client/cli/genesis_accts.go
func addGenesisAccountCmd(ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-account [address_or_key_name] [coin][,[coin]]",
		Short: "Add genesis account to genesis.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				kb, err := keys.NewKeyBaseFromDir(viper.GetString(flagClientHome))
				if err != nil {
					return err
				}

				info, err := kb.Get(args[0])
				if err != nil {
					return err
				}

				addr = info.GetAddress()
			}

			coins, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			genAcc := genaccounts.NewGenesisAccountRaw(addr, coins, sdk.NewCoins(), 0, 0, "")
			if err := genAcc.Validate(); err != nil {
				return err
			}

			// retrieve the app state
			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return err
			}

			// add genesis account to the app state
			var genesisAccounts genaccounts.GenesisAccounts

			if appState[genaccounts.ModuleName] != nil {
				cdc.MustUnmarshalJSON(appState[genaccounts.ModuleName], &genesisAccounts)
			}

			if genesisAccounts.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %v", addr)
			}

			genesisAccounts = append(genesisAccounts, genAcc)

			genesisStateBz := cdc.MustMarshalJSON(genaccounts.GenesisState(genesisAccounts))
			appState[genaccounts.ModuleName] = genesisStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			// export app state
			genDoc.AppState = appStateJSON

			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().String(flagClientHome, defaultClientHome, "client's home directory")
	return cmd
}
