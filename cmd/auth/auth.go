package auth

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/z411392/golab/adapters/http/dummy_json"
)

const use = "auth"

var (
	username string
	password string
)

func runE(command *cobra.Command, args []string) (err error) {
	adapter := dummy_json.NewDummyJsonAdapter()
	credentials, err := adapter.SignIn(username, password)
	if err != nil {
		return
	}
	if credentials == nil || credentials.AccessToken == "" {
		err = fmt.Errorf("authentication failed")
	}
	fmt.Printf("%s\n", credentials.AccessToken)
	return
}

func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:  use,
		RunE: runE,
	}
	flags := command.Flags()
	flags.StringVarP(&username, "username", "u", "", "")
	flags.StringVarP(&password, "password", "p", "", "")
	return command
}
