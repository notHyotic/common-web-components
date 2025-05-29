package main

import (
	"lesiw.io/ops"
	commands "ops/commands"
)

func main() {
	ops.Handle(commands.Ops{})
}