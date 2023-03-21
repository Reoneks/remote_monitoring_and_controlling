package main

import (
	"project/cmd"

	"go.uber.org/fx"
)

func main() {
	fx.New(cmd.Exec()).Run()
}
