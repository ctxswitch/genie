package main

import (
	"ctx.sh/genie/pkg/cmd"
)

func main() {
	root := cmd.NewRoot()
	root.Execute()
}
