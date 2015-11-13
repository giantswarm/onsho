package main

import (
	"fmt"

	"github.com/giantswarm/moa/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type DestroyFlags struct {
	TMuxSessionName string
}

var (
	destroyCmd = &cobra.Command{
		Use:   "destroy [serial|all]",
		Short: "Destroy virtual machines.",
		Long:  "Destroy virtual machines.",
		Run:   destroyRun,
	}

	destroyFlags = &DestroyFlags{}
)

func init() {
	destroyCmd.PersistentFlags().StringVar(&destroyFlags.TMuxSessionName, "tmux-session-name", "zoo", "TMUX session name to start the instances in")
}

func destroyRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Serial missing. If you want to destroy all vms pass in 'all'.")
	}

	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	if args[0] == "all" {
		vm.StopAll(destroyFlags.TMuxSessionName)
		vm.DestroyAll(configDir)
	} else {
		machine, err := vm.Load(configDir, args[0])
		assert(err)

		machine.Stop()
		machine.Destroy()
	}
}
