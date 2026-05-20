package main

import (
	"github.com/scottbass3/altcha-server/internal/command"
)

var (
	ProjectVersion = "dev"
	GitRef         = "unknown"
	BuildDate      = "unknown"
)

func main() {
	command.Main(
		ProjectVersion+" ("+GitRef+") "+BuildDate,
		command.RunCommand(),
		command.GenerateCommand(),
		command.SolveCommand(),
		command.VerifyCommand(),
	)
}
