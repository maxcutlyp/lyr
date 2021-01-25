package cmd

import (
	"github.com/satoqz/lyr/config"
	"github.com/spf13/cobra"
)

var tokenCmd = &cobra.Command{
	Use:  "token",
	Args: cobra.MinimumNArgs(1),
	RunE: tokenExec,
}

func tokenExec(_ *cobra.Command, args []string) error {
	return config.SetToken(args[0])
}
