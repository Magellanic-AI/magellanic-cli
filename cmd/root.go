// Package cmd /*
// Copyright © 2023 Magellanic <contact@magellanic.ai>
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

const ViperPrefix = "mgl_"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: "magellanic-cli",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		dotenvPath, err := cmd.Flags().GetString("cli_config_path")
		if err != nil {
			return err
		}
		if dotenvPath != "" {
			viper.SetConfigType("env")
			viper.SetConfigFile(dotenvPath)
			_ = viper.ReadInConfig()
		}
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			viper.BindPFlag(ViperPrefix+f.Name, f)
			if !f.Changed && viper.IsSet(ViperPrefix+f.Name) {
				val := viper.Get(ViperPrefix + f.Name)
				cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
			}
		})
		err = validateParams(cmd)
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("api_key", "a", "", "API key")
	viper.BindEnv("mgl_api_key", "MGL_API_KEY")
	RootCmd.PersistentFlags().StringP("api_url", "u", "https://api.magellanic.ai", "API url")
	viper.BindEnv("mgl_api_url", "MGL_API_URL")
	RootCmd.PersistentFlags().StringP("cli_config_path", "C", "~/.magellanic/creds", "configuration file path")
	viper.BindEnv("mgl_cli_config_path", "MGL_CLI_CONFIG_PATH")
}

func validateParams(cmd *cobra.Command) error {
	apiKey, err := cmd.Flags().GetString("api_key")
	if err != nil {
		return err
	}
	if apiKey == "" {
		return errors.New("API key not defined")
	}
	return nil
}
