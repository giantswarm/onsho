package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/giantswarm/moa/tmux"
	"github.com/giantswarm/moa/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create virtual machines.",
		Long:  "Create virtual machines.",
		Run:   createRun,
	}

	createFlags = &vm.VMFlags{}
)

func init() {
	createCmd.PersistentFlags().Uint8Var(&createFlags.NumberOfVMs, "num-vms", 5, "Number of virtual machines")
	createCmd.PersistentFlags().StringVar(&createFlags.BridgeInterfaces, "bridge-ifs", "bond0", "Bridge interface to bind the virtual machines to (comma-separated)")
	createCmd.PersistentFlags().StringVar(&createFlags.DiskSize, "disk-size", "16G", "Disk size of the virtual machines in GB (eg 16G)")
	createCmd.PersistentFlags().BoolVar(&createFlags.NoTmux, "no-tmux", false, "Run a single vm within the current shell")
	createCmd.PersistentFlags().StringVar(&createFlags.TMuxSessionName, "tmux-session-name", "zoo", "TMUX session name to start the instances in")
	createCmd.PersistentFlags().StringVar(&createFlags.HDs, "hds", "hda,hdb", "Names of the hard disk devices (comma-separated)")
	createCmd.PersistentFlags().Uint16Var(&createFlags.Memory, "memory", 1024, "RAM of the virtual machines in MB (eg 1024)")
	createCmd.PersistentFlags().StringVar(&createFlags.Image, "image", "ipxe.iso", "Image to start the virtual machine with")
}

func createRun(cmd *cobra.Command, args []string) {
	if createFlags.NoTmux && createFlags.NumberOfVMs > 1 {
		fmt.Println("You can only start a single VM without tmux!")
		os.Exit(1)
	}

	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	if !createFlags.NoTmux {
		tmux.KillSession(createFlags.TMuxSessionName)
	}

	// create vms
	for i := 0; i < int(createFlags.NumberOfVMs); i++ {
		machine, err := vm.NewVM(createFlags, configDir)
		assert(err)

		qemuArgs := machine.GenerateQEMUArgs()
		qemuCmd := fmt.Sprintf("%s %s", "qemu-system-x86_64", strings.Join(qemuArgs, " "))
		if globalFlags.debug || globalFlags.verbose {
			fmt.Println(qemuCmd)
		}

		if !createFlags.NoTmux {
			tmux.NewWindow(createFlags.TMuxSessionName, qemuCmd)
		} else {
			cmd := exec.Command("qemu-system-x86_64", qemuArgs...)

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("%s %s - %v\n", stdout.String(), stderr.String(), err)
				os.Exit(1)
			}
		}
	}
}
