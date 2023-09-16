/*
Copyright © 2023 Albert David Lewandowski a.lewandowski@magellanic.ai
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var (
	RootCmd = &cobra.Command{
		Use: "magellanic-cli",
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			dotenvPath, err := cmd.Flags().GetString("dotenv_path")
			if err != nil {
				return err
			}
			if dotenvPath != "" {
				viper.SetConfigType("env")
				viper.SetConfigFile(dotenvPath)
				if err := viper.ReadInConfig(); err != nil {
					return err
				}
			}
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				viper.BindPFlag(f.Name, f)
				if !f.Changed && viper.IsSet(f.Name) {
					val := viper.Get(f.Name)
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
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringP("api_key", "a", "", "API key")
	viper.BindEnv("api_key", "API_KEY")
	RootCmd.PersistentFlags().StringP("api_url", "u", "https://api.magellanic.ai", "API url")
	viper.BindEnv("api_url", "API_URL")
	RootCmd.PersistentFlags().StringP("dotenv_path", "d", "", ".env configuration file path")
	viper.BindEnv("dotenv_path", "DOTENV_PATH")
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
