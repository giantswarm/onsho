package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/giantswarm/onsho/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type CreateFlags struct {
	NumberOfVMs     uint8
	NoStart         bool
	Population      string
	NoTMux          bool
	TMuxSessionName string
}

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create virtual machines.",
		Long:  "Create virtual machines.",
		Run:   createRun,
	}

	vmFlags     = &vm.VMFlags{}
	createFlags = &CreateFlags{}
)

func init() {
	createCmd.PersistentFlags().Uint8Var(&createFlags.NumberOfVMs, "num-vms", DefaultNumVMs, "Number of virtual machines")
	createCmd.PersistentFlags().BoolVar(&createFlags.NoStart, "no-start", false, "Do not start the virtual machines")
	createCmd.PersistentFlags().StringVar(&createFlags.Population, "population", "", "List of serials and macs to populate the machines from")
	createCmd.PersistentFlags().BoolVar(&createFlags.NoTMux, "no-tmux", false, "Run a single vm within the current shell")
	createCmd.PersistentFlags().StringVar(&createFlags.TMuxSessionName, "tmux-session-name", DefaultTMuxSessionName, "TMUX session name to start the instances in")

	createCmd.PersistentFlags().StringVar(&vmFlags.BridgeInterfaces, "bridge-ifs", DefaultBridgeIfs, "Bridge interface to bind the virtual machines to (comma-separated)")
	createCmd.PersistentFlags().StringVar(&vmFlags.DiskSize, "disk-size", DefaultDiskSize, "Disk size of the virtual machines in GB (eg 16G)")
	createCmd.PersistentFlags().StringVar(&vmFlags.HDs, "hds", DefaultHDs, "Names of the hard disk devices (comma-separated)")
	createCmd.PersistentFlags().Uint16Var(&vmFlags.Memory, "memory", DefaultMemory, "RAM of the virtual machines in MB (eg 1024)")
	createCmd.PersistentFlags().StringVar(&vmFlags.Image, "image", DefaultImage, "Image to start the virtual machine with")
}

func loadPopulation(filePath string) (vm.Population, error) {
	pop := vm.Population{}

	f, err := os.Open(filePath)
	if err != nil {
		return pop, err
	}
	defer f.Close()

	popBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return pop, err
	}

	err = yaml.Unmarshal(popBytes, &pop)
	return pop, err

}

func createRun(cmd *cobra.Command, args []string) {
	if createFlags.NoTMux && createFlags.NumberOfVMs > 1 {
		fmt.Println("You can only start a single VM without tmux!")
		os.Exit(1)
	}

	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	pop := vm.Population{}
	if createFlags.Population != "" {
		pop, err = loadPopulation(createFlags.Population)
		assert(err)
	}

	for i := 0; i < int(createFlags.NumberOfVMs); i++ {
		machine, err := vm.NewVM(vmFlags, configDir, &pop)
		assert(err)

		if !createFlags.NoStart {
			machine.Start(createFlags.TMuxSessionName, createFlags.NoTMux)
		}
	}
}
