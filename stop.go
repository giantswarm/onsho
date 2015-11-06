package main

import (
	"github.com/spf13/cobra"
)

var (
	stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop virtual machines.",
		Long:  "Stop virtual machines.",
		Run:   stopRun,
	}
)

func stopRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}
