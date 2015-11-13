package main

import (
	"fmt"

	"github.com/giantswarm/moa/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type StartFlags struct {
	NoTMux          bool
	TMuxSessionName string
}

var (
	startCmd = &cobra.Command{
		Use:   "start [serial]",
		Short: "Start virtual machines.",
		Long:  "Start virtual machines.",
		Run:   startRun,
	}

	startFlags = &StartFlags{}
)

func init() {
	startCmd.PersistentFlags().BoolVar(&startFlags.NoTMux, "no-tmux", false, "Run a single vm within the current shell")
	startCmd.PersistentFlags().StringVar(&startFlags.TMuxSessionName, "tmux-session-name", "zoo", "TMUX session name to start the instances in")
}

func startRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Serial missing")
	}

	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	machine, err := vm.Load(configDir, args[0])
	assert(err)

	machine.Start(startFlags.TMuxSessionName, startFlags.NoTMux)
}
