package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/giantswarm/onsho/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

type StatusFlags struct {
	TMuxSessionName string
}

var (
	statusCmd = &cobra.Command{
		Use:   "status [serial]",
		Short: "Status virtual machines.",
		Long:  "Status virtual machines.",
		Run:   statusRun,
	}

	statusFlags = &StatusFlags{}
)

const (
	statusScheme = "%s | %s"
)

func init() {
	statusCmd.PersistentFlags().StringVar(&statusFlags.TMuxSessionName, "tmux-session-name", DefaultTMuxSessionName, "TMUX session name to status the instances in")
}

func statusRun(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("Serial missing")
		os.Exit(1)
	}

	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	machine, err := vm.Load(configDir, args[0])
	assert(err)

	status, err := machine.Status(statusFlags.TMuxSessionName)
	assert(err)

	var hds, nics []string
	for _, hd := range machine.HDs {
		hds = append(hds, fmt.Sprintf("%s %s", hd.Device, hd.Size))
	}
	for _, nic := range machine.NICs {
		nics = append(nics, fmt.Sprintf("%s %s", nic.Bridge, nic.Mac))
	}

	lines := []string{}
	lines = append(lines, fmt.Sprintf(statusScheme, "Serial:", machine.Serial))
	lines = append(lines, fmt.Sprintf(statusScheme, "Status:", status))
	lines = append(lines, fmt.Sprintf(statusScheme, "Image:", machine.Image))
	lines = append(lines, fmt.Sprintf("%s | %d", "Memory:", machine.Memory))
	lines = append(lines, fmt.Sprintf(statusScheme, "HDs:", strings.Join(hds, ",")))
	lines = append(lines, fmt.Sprintf(statusScheme, "NICs:", strings.Join(nics, ",")))
	fmt.Println(columnize.SimpleFormat(lines))
}
