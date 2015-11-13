package main

import (
	"fmt"

	"github.com/giantswarm/moa/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type StopFlags struct {
	TMuxSessionName string
}

var (
	stopCmd = &cobra.Command{
		Use:   "stop [serial|all]",
		Short: "Stop virtual machines.",
		Long:  "Stop virtual machines.",
		Run:   stopRun,
	}

	stopFlags = &StopFlags{}
)

func init() {
	stopCmd.PersistentFlags().StringVar(&stopFlags.TMuxSessionName, "tmux-session-name", "zoo", "TMUX session name to stop the instances in")
}

func stopRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Serial missing. If you want to stop all vms pass in 'all'.")
	}

	if args[0] == "all" {
		vm.StopAll(stopFlags.TMuxSessionName)
	} else {
		configDir, err := homedir.Expand(globalFlags.config)
		assert(err)

		machine, err := vm.Load(configDir, args[0])
		assert(err)

		machine.Stop()
	}
}
