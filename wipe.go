package main

import (
	"fmt"
	"os"

	"github.com/giantswarm/moa/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	wipeCmd = &cobra.Command{
		Use:   "wipe",
		Short: "Wipe a virtual machine.",
		Long:  "Wipe a virtual machine and recreate the disks.",
		Run:   wipeRun,
	}

	wipeVMFlags = &vm.VMFlags{}
)

func init() {
	wipeCmd.PersistentFlags().StringVar(&wipeVMFlags.DiskSize, "disk-size", DefaultDiskSize, "Disk size of the virtual machines in GB (eg 16G)")
	wipeCmd.PersistentFlags().StringVar(&wipeVMFlags.HDs, "hds", DefaultHDs, "Names of the hard disk devices (comma-separated)")
}

func wipeRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Serial missing")
		os.Exit(1)
	}

	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	machine, err := vm.Load(configDir, args[0])
	assert(err)

	machine.Wipe(wipeVMFlags)

	err = machine.Save()
	assert(err)
}
