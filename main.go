package main

import (
	"remote_monitoring_and_controlling/cmd"

	"go.uber.org/fx"
)

func main() {
	fx.New(cmd.Exec()).Run()
}
