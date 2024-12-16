package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/z411392/golab/cmd/auth"
	"github.com/z411392/golab/cmd/cdc"
	"github.com/z411392/golab/cmd/serve"
	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	maxprocs.Set()
	command := &cobra.Command{
		Use: os.Getenv("FIREBASE_PROJECT_ID"),
	}
	command.AddCommand(serve.NewCommand())
	command.AddCommand(auth.NewCommand())
	command.AddCommand(cdc.NewCommand())
	command.Execute()
}
