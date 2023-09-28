// Package config /*
// Copyright © 2023 Magellanic <contact@magellanic.ai>
package config

import (
	"magellanic-cli/cmd"

	"github.com/spf13/cobra"
)

const apiPath = "/public-api/configs"

var configCmd = &cobra.Command{
	Use: "config",
}

func init() {
	cmd.RootCmd.AddCommand(configCmd)
}
