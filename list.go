package main

import (
	"fmt"
	"strings"

	"github.com/giantswarm/moa/vm"
	"github.com/mitchellh/go-homedir"
	"github.com/ryanuber/columnize"
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

const (
	listHeader = "Id | Image | Memory | Disks | NICs"
	listScheme = "%s | %s | %d | %s | %s"
)

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

		lines = append(lines, fmt.Sprintf(listScheme,
			machine.Serial,
			machine.Image,
			machine.Memory,
			strings.Join(hds, ","),
			strings.Join(nics, ",")))
	}
	fmt.Println(columnize.SimpleFormat(lines))
}
