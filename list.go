package main

import (
	"fmt"
	"strings"

	"github.com/giantswarm/onsho/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/ryanuber/columnize"
	"github.com/spf13/cobra"
)

type ListFlags struct {
	TMuxSessionName string
}

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List virtual machines.",
		Long:  "List virtual machines.",
		Run:   listRun,
	}

	listFlags = &ListFlags{}
)

const (
	listHeader = "Id | Image | Memory | Disks | NICs | Status"
	listScheme = "%s | %s | %d | %s | %s | %s"
)

func init() {
	listCmd.PersistentFlags().StringVar(&listFlags.TMuxSessionName, "tmux-session-name", DefaultTMuxSessionName, "TMUX session name to status the instances in")
}

func listRun(cmd *cobra.Command, args []string) {
	configDir, err := homedir.Expand(globalFlags.config)
	assert(err)

	machines, err := vm.List(configDir)
	assert(err)

	lines := []string{listHeader}
	for _, machine := range machines {
		var hds, nics []string
		for _, hd := range machine.HDs {
			hds = append(hds, fmt.Sprintf("%s %s", hd.Device, hd.Size))
		}
		for _, nic := range machine.NICs {
			nics = append(nics, fmt.Sprintf("%s %s", nic.Bridge, nic.Mac))
		}

		status, _ := machine.Status(listFlags.TMuxSessionName)

		lines = append(lines, fmt.Sprintf(listScheme,
			machine.Serial,
			machine.Image,
			machine.Memory,
			strings.Join(hds, ","),
			strings.Join(nics, ","),
			status))
	}
	fmt.Println(columnize.SimpleFormat(lines))
}
