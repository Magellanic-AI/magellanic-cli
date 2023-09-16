/*
Copyright © 2023 Albert David Lewandowski a.lewandowski@magellanic.ai
*/
package config

import (
	"magellanic-cli/cmd"

	"github.com/spf13/cobra"
)

const apiPath = "/public-api/configs"

// configCmd represents the config/config command
var configCmd = &cobra.Command{
	Use: "config",
}

func init() {
	cmd.RootCmd.AddCommand(configCmd)
}
