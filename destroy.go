package main

import (
	"github.com/spf13/cobra"
)

var (
	destroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy virtual machines.",
		Long:  "Destroy virtual machines.",
		Run:   destroyRun,
	}
)

func destroyRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}
