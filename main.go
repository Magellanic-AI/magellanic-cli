/*
Copyright © 2023 Magellanic <contact@magellanic.ai>
*/
package main

import (
	"magellanic-cli/cmd"
	_ "magellanic-cli/cmd/config"
)

func main() {
	cmd.Execute()
}
