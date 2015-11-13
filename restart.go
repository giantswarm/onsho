package main

import (
	"github.com/spf13/cobra"
)

var (
	restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart virtual machines.",
		Long:  "Restart virtual machines.",
		Run:   restartRun,
	}
)

func restartRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}
