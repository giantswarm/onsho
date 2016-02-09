package main

import (
	"fmt"
	"os"

	"github.com/giantswarm/onsho/vm"
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
	stopCmd.PersistentFlags().StringVar(&stopFlags.TMuxSessionName, "tmux-session-name", DefaultTMuxSessionName, "TMUX session name to stop the instances in")
}

func stopRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Serial missing. If you want to stop all vms pass in 'all'.")
		os.Exit(1)
	}

	if args[0] == "all" {
		vm.StopAll(stopFlags.TMuxSessionName)
	} else {
		configDir, err := homedir.Expand(globalFlags.config)
		assert(err)

		machine, err := vm.Load(configDir, args[0])
		assert(err)

		err = machine.Stop(stopFlags.TMuxSessionName)
		assert(err)
	}
}
