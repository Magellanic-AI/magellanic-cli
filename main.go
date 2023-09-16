/*
Copyright © 2023 Albert David Lewandowski a.lewandowski@magellanic.ai
*/
package main

import (
	"magellanic-cli/cmd"
	_ "magellanic-cli/cmd/config"
)

func main() {
	cmd.Execute()
}
