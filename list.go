package main

import (
	"fmt"

	"github.com/giantswarm/moa/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List virtual machines.",
		Long:  "List virtual machines.",
		Run:   listRun,
	}
)

func listRun(cmd *cobra.Command, args []string) {
	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	list, err := vm.List(configDir)
	assert(err)

	for _, machine := range list {
		fmt.Println(machine.Serial)
	}
}
