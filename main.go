package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	DefaultNumVMs          uint8  = 5
	DefaultTMuxSessionName string = "zoo"
	DefaultBridgeIfs       string = "bond0"
	DefaultDiskSize        string = "16G"
	DefaultHDs             string = "hda,hdb"
	DefaultMemory          uint16 = 1024
	DefaultImage           string = "ipxe.iso"
)

var (
	globalFlags struct {
		debug   bool
		verbose bool
		config  string
		sleep   time.Duration
	}

	mainCmd = &cobra.Command{
		Use:   "moa",
		Short: "Manage a QEMU Giant Swarm",
		Long:  "Manage a QEMU Giant Swarm",
		Run:   mainRun,
	}

	projectVersion string
	projectBuild   string
)

func init() {
	mainCmd.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "d", false, "Print debug output")
	mainCmd.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Print verbose output")
	mainCmd.PersistentFlags().DurationVar(&globalFlags.sleep, "sleep", 60*time.Second, "time to sleep between machine starts")
	mainCmd.PersistentFlags().StringVar(&globalFlags.config, "config", "~/.giantswarm/moa", "Config folder (for machine state and boot images)")
}

func createConfig() {
	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0700)
	}
	machinesDir := configDir + "/machines"
	if _, err := os.Stat(machinesDir); os.IsNotExist(err) {
		os.MkdirAll(machinesDir, 0700)
	}
	disksDir := configDir + "/disks"
	if _, err := os.Stat(disksDir); os.IsNotExist(err) {
		os.MkdirAll(disksDir, 0700)
	}
}

func assert(err error) {
	if err != nil {
		if globalFlags.debug {
			fmt.Printf("%#v\n", err)
			os.Exit(1)
		} else {
			log.Fatal(err)
		}
	}
}

func confirm(question string) error {
	for {
		fmt.Printf("%s ", question)
		bio := bufio.NewReader(os.Stdin)
		line, _, err := bio.ReadLine()
		if err != nil {
			return err
		}

		if string(line) == "yes" {
			return nil
		}
		fmt.Println("Please enter 'yes' to confirm.")
	}
}

func mainRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func main() {
	createConfig()
	mainCmd.AddCommand(versionCmd)
	mainCmd.AddCommand(createCmd)
	mainCmd.AddCommand(destroyCmd)
	mainCmd.AddCommand(stopCmd)
	mainCmd.AddCommand(startCmd)
	mainCmd.AddCommand(restartCmd)
	mainCmd.AddCommand(listCmd)
	mainCmd.AddCommand(wipeCmd)

	mainCmd.Execute()
}
