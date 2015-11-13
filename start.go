package main

import (
	"github.com/spf13/cobra"
)

var (
	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start virtual machines.",
		Long:  "Start virtual machines.",
		Run:   startRun,
	}
)

func startRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}
