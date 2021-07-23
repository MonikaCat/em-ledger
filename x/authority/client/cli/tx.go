// This software is Copyright (c) 2019-2020 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package cli

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgtypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/e-money/em-ledger/util"
	"github.com/e-money/em-ledger/x/authority/types"
	"github.com/spf13/cobra"
)

func GetTxCmd() *cobra.Command {
	authorityCmds := &cobra.Command{
		Use:                "authority",
		Short:              "Manage authority tasks",
		DisableFlagParsing: false,
	}

	authorityCmds.AddCommand(
		getCmdCreateIssuer(),
		getCmdDestroyIssuer(),
		getCmdSetGasPrices(),
		GetCmdReplaceAuthority(),
		GetCmdScheduleUpgrade(),
	)

	return authorityCmds
}

func getCmdSetGasPrices() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-gas-prices [authority_key_or_address] [minimum_gas_prices]",
		Example: "emd tx authority set-gas-prices masterkey 0.0005eeur,0.0000001ejpy",
		Short:   "Control the minimum gas prices for the chain",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasPrices, err := sdk.ParseDecCoins(args[1])
			if err != nil {
				return err
			}

			msg := &types.MsgSetGasPrices{
				GasPrices: gasPrices,
				Authority: clientCtx.GetFromAddress().String(),
			}
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdCreateIssuer() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-issuer [authority_key_or_address] [issuer_address] [denominations]",
		Example: "emd tx authority create-issuer masterkey emoney17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu eeur,ejpy",
		Short:   "Create a new issuer",
		Args:    cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			issuerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			denoms, err := util.ParseDenominations(args[2])
			if err != nil {
				return err
			}

			msg := &types.MsgCreateIssuer{
				Issuer:        issuerAddr.String(),
				Denominations: denoms,
				Authority:     clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func getCmdDestroyIssuer() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "destroy-issuer [authority_key_or_address] [issuer_address]",
		Example: "emd tx authority destory-issuer masterkey emoney17up20gamd0vh6g9ne0uh67hx8xhyfrv2lyazgu",
		Short:   "Delete an issuer",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Flags().Set(flags.FlagFrom, args[0])
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			issuerAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := &types.MsgDestroyIssuer{
				Issuer:    issuerAddr.String(),
				Authority: clientCtx.GetFromAddress().String(),
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdReplaceAuthority() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "replace [authority_key_or_address] new_authority_address",
		Short:   "Replace the authority key",
		Example: "emd tx authority replace emoney1n5ggspeff4fxc87dvmg0ematr3qzw5l4v20mdv emoney1hq6tnhqg4t7358f3vd9crru93lv0cgekdxrtgv",
		Long: `Replace the authority key with a new multisig address. 
For a 24-hour grace period the former authority key is equivalent to the new one.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			f := cmd.Flags()

			err := f.Set(flags.FlagFrom, args[0])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgReplaceAuthority{
				Authority:    clientCtx.GetFromAddress().String(),
				NewAuthority: args[1],
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func GetCmdScheduleUpgrade() *cobra.Command {
	const (
		upgHeight = "upg-height"
		upgTime   = "upg-time"
		upgInfo   = "upg-info"
	)

	var (
		upgHeightVal, upgTimeSecsVal int64
		upgInfoVal                   string
	)

	cmd := &cobra.Command{
		Use:   "schedule-upg [authority_key_or_address] plan_name",
		Short: "Schedule a software upgrade.",
		Example: `emd tx schedule-upg emoney1n5ggspeff4fxc87dvmg0ematr3qzw5l4v20mdv 0.43 --upg_height 2001
emd tx schedule-upg emoney1n5ggspeff4fxc87dvmg0ematr3qzw5l4v20mdv 'New Staking Rewards 36%' --upg_time 1628956125 # Unix seconds for 2021-08-14 15:48:45 +0000 UTC
emd tx schedule-upg emoney1n5ggspeff4fxc87dvmg0ematr3qzw5l4v20mdv sdk-v0.43.0 --upg_height 2001 --upg_info "https://e-money.com/mainnet-099-info.json?checksum=sha256:deaaa99fda9407c4dbe1d04bd49bab0cc3c1dd76fa392cd55a9425be074af01e"`,
		Long: `Schedule a software upgrade.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := cmd.Flags().Set(flags.FlagFrom, args[0])
			if err != nil {
				return err
			}

			if err := validateUpgFlags(
				upgHeight, upgHeightVal, upgTime, upgTimeSecsVal,
			); err != nil {
				return err
			}

			upgTimeVal := time.Unix(upgTimeSecsVal, 0)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := &types.MsgScheduleUpgrade{
				Authority: clientCtx.GetFromAddress().String(),
				Plan: upgtypes.Plan{
					Name:   args[1],
					Time:   upgTimeVal,
					Height: upgHeightVal,
					Info:   upgInfoVal,
				},
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	f := cmd.Flags()
	f.Int64VarP(&upgHeightVal, upgHeight, "n", 0, "upgrade block height number")
	f.Int64VarP(
		&upgTimeSecsVal, upgTime, "t", 0, "upgrade block time (in Unix seconds)",
	)
	f.StringVarP(
		&upgInfoVal, upgInfo, "i", "", "upgrade info",
	)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func validateUpgFlags(
	upgHeight string, upgHeightVal int64, upgTime string, timeVal int64,
) error {
	if upgHeightVal == 0 && timeVal == 0 {
		return sdkerrors.Wrapf(
			types.ErrMissingFlag,
			"need to specify --%s or --%s", upgHeight, upgTime,
		)
	}

	flagsSet := 0
	if upgHeightVal != 0 {
		flagsSet++
	}
	if timeVal != 0 {
		flagsSet++
	}
	if flagsSet != 1 {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest,
			"specify only one of the flags: --%s or --%s", upgHeight,
			upgTime,
		)
	}

	return nil
}
