package main

import (
	"github.com/scottbass3/altcha-server/internal/command"
)

func main() {
	command.Main(
		command.RunCommand(),
		command.GenerateCommand(),
		command.SolveCommand(),
		command.VerifyCommand(),
	)
}
