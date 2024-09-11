package main

import (
	"forge.cadoles.com/cadoles/altcha-server/internal/command"
)

func main() {
	command.Main(
		command.RunCommand(),
		command.GenerateCommand(),
		command.SolveCommand(),
		command.VerifyCommand(),
	)
}