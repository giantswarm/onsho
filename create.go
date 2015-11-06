package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/giantswarm/moa/tmux"
	"github.com/spf13/cobra"
)

type CreateFlags struct {
	NumberOfVMs      uint8
	BridgeInterfaces string
	DiskSize         string
	NoTmux           bool
	TMuxSessionName  string
	HDs              string
	Memory           uint16
	Image            string
}

type VM struct {
	hds    []HD
	serial string
	nics   []NIC
}

type HD struct {
	device string
	path   string
}

type NIC struct {
	bridge string
	mac    string
	dev    string
	addr   string
}

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create virtual machines.",
		Long:  "Create virtual machines.",
		Run:   createRun,
	}

	createFlags = &CreateFlags{}

	uuids = []string{
		"004b27ed-692e-b32e-1f68-d89aff66c71b",
		"aa1f18e1-f14f-2dd9-4fa0-dae7317c712c",
		"7100c054-d2c9-e299-b669-e8bdb85f6904",
		"2843c49e-d1ba-6dd3-1320-d7cc82d8ea3a",
		"175842d1-ce55-da90-ab2a-308b532aa17b",
		"ed136c8b-ad6d-d604-89f3-29262a63fc76",
		"bcd0b26c-33c8-7bb7-370d-0200e1246f61",
		"b9eae062-101a-b17e-3728-94a3f85ccb74",
		"4acb7094-30e9-fa14-8cfc-6403b065177b",
		"d523b3fd-fe7b-757b-2610-7a88433e7e6a",
		"1b7e3615-68ba-19b5-4050-9e274e96c933",
		"a03b788b-fc8c-0684-96df-38b81daf77a3",
		"0ecdc715-3db7-3ef7-8491-7ea93be48d60",
		"98d53981-dde4-b2a0-9a99-db5f0e7674e8",
		"19cab58a-721a-7771-d140-95e0ef559f30",
		"135218ab-2056-36b8-8a64-106cff459fa8",
		"86bc2477-80b7-3d29-3290-f0c6c7c899f8",
		"ce20c3c0-f69a-2eeb-5a43-6076e614f699",
		"5f4056fa-2fc5-bbb7-9b3e-a7dba2b00b3e",
		"7a14db43-ee98-5d3c-6e9d-5e324d529eaa",
	}

	macs = []string{
		"00:16:3e:a0:b7:df",
		"00:16:3e:ca:0a:19",
		"00:16:3e:c5:b0:5b",
		"00:16:3e:fb:16:9a",
		"00:16:3e:1f:9b:3d",
		"00:16:3e:bf:f6:06",
		"00:16:3e:87:ad:cf",
		"00:16:3e:98:20:39",
		"00:16:3e:09:6a:2b",
		"00:16:3e:a0:02:5c",
		"00:16:3e:9e:5e:0d",
		"00:16:3e:7c:36:2d",
		"00:16:3e:ce:f2:d9",
		"00:16:3e:a7:d0:d5",
		"00:16:3e:f6:18:8d",
		"00:16:3e:1d:a0:e8",
		"00:16:3e:0c:94:1d",
		"00:16:3e:d6:44:5e",
		"00:16:3e:6d:85:28",
		"00:16:3e:17:3b:f3",
		"00:16:3e:a7:6d:38",
		"00:16:3e:c4:a9:8d",
		"00:16:3e:84:95:d4",
		"00:16:3e:d3:aa:5c",
		"00:16:3e:44:f5:58",
		"00:16:3e:72:eb:e8",
		"00:16:3e:bb:58:59",
		"00:16:3e:63:54:5a",
		"00:16:3e:6c:0b:74",
		"00:16:3e:74:8e:8e",
		"00:16:3e:b9:73:7c",
		"00:16:3e:3f:3d:be",
		"00:16:3e:5a:2d:4b",
		"00:16:3e:56:b0:01",
		"00:16:3e:63:c1:4d",
		"00:16:3e:22:7c:4a",
		"00:16:3e:26:6d:30",
		"00:16:3e:c8:c2:a1",
		"00:16:3e:4a:a6:7f",
		"00:16:3e:db:43:dd",
		"00:16:3e:b6:34:03",
		"00:16:3e:39:9a:3f",
		"00:16:3e:90:1e:93",
		"00:16:3e:51:33:1b",
		"00:16:3e:a0:b6:7e",
		"00:16:3e:9a:34:28",
		"00:16:3e:9a:5b:82",
		"00:16:3e:ad:72:f5",
		"00:16:3e:4a:70:9b",
		"00:16:3e:93:49:95",
		"00:16:3e:d7:81:3e",
		"00:16:3e:5e:de:69",
		"00:16:3e:3e:17:f0",
		"00:16:3e:5d:10:d8",
		"00:16:3e:83:5e:1c",
		"00:16:3e:7d:47:f3",
		"00:16:3e:d0:f7:07",
		"00:16:3e:ea:a5:41",
		"00:16:3e:9c:8f:b7",
		"00:16:3e:f4:f3:ef",
	}

	networkInterfaceAddresses = []string{
		"03", "05", "07",
	}
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

func validateBridges(bridges []string) bool {
	validateOk := true
	// validate has interface (iterate through bridge interfaces)
	for _, b := range bridges {
		if _, err := os.Stat("/sys/class/net/" + b); os.IsNotExist(err) {
			validateOk = false
			fmt.Printf(`Interface %s doesn't exist, you might want to set it up:

sysctl net.ipv4.ip_forward=1

brctl addbr %s
ip add add cidr_addr dev %s
ip link set %s up
`, b, b, b, b)
		}
	}
	return validateOk
}

func validateNumberOfVMs(numberOfVMs uint8, bridges []string) bool {
	validateOk := true
	if len(uuids) < int(numberOfVMs) {
		validateOk = false
		fmt.Printf("dude... no. You need to feed me more random uuids than %d.\n", numberOfVMs)
	}
	if len(macs) < (len(bridges) * int(numberOfVMs)) {
		validateOk = false
		fmt.Printf("dude... no. You need to feed me more random mac addresses than %d.\n", len(bridges)*int(numberOfVMs))
	}
	return validateOk
}

func createVM(number int, hds []string, bridges []string, baseDir string) VM {
	vm := VM{
		serial: uuids[number],
	}

	for _, hd := range hds {
		vm.hds = append(vm.hds, HD{
			device: hd,
			path:   fmt.Sprintf("%s/%s-%s.qcow2", baseDir, vm.serial, hd),
		})
	}

	for i, bridge := range bridges {
		vm.nics = append(vm.nics, NIC{
			bridge: bridge,
			mac:    macs[0],
			dev:    "e1000",
			addr:   networkInterfaceAddresses[i],
		})
		macs = macs[1:]
	}

	return vm
}

func ensureHDExists(path, size string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if globalFlags.debug || globalFlags.verbose {
			fmt.Println("qemu-img", "create", "-f", "qcow2", path, size)
		}
		cmd := exec.Command("qemu-img", "create", "-f", "qcow2", path, size)

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

func generateQEMUArgs(vm VM, memory uint16, image string) []string {
	var args []string
	args = append(args, "-m", strconv.Itoa(int(memory)))
	args = append(args, "-cdrom", image)
	args = append(args, "-boot", "d")
	args = append(args, "-nographic")
	args = append(args, "-device", "sga")
	args = append(args, "-serial", "mon:stdio")

	for vlan, nic := range vm.nics {
		args = append(args, "-net", fmt.Sprintf("bridge,br=%s,vlan=%d", nic.bridge, vlan))
		args = append(args, "-net", fmt.Sprintf("nic,vlan=%d,model=%s,macaddr=%s,addr=%s", vlan, nic.dev, nic.mac, nic.addr))
	}

	args = append(args, "--cpu", "host")
	args = append(args, "-enable-kvm")

	for _, hd := range vm.hds {
		args = append(args, fmt.Sprintf("-%s", hd.device), hd.path)
	}

	args = append(args, "-smbios", fmt.Sprintf("type=3,serial=%s", vm.serial))
	args = append(args, "-monitor", fmt.Sprintf("unix:/tmp/tinyswarm-%s.sock,server,nowait", vm.serial))

	return args
}

func createRun(cmd *cobra.Command, args []string) {
	if createFlags.NoTmux && createFlags.NumberOfVMs > 1 {
		fmt.Println("You can only start a single VM without tmux!")
		os.Exit(1)
	}

	hds := strings.Split(createFlags.HDs, ",")
	bridges := strings.Split(createFlags.BridgeInterfaces, ",")

	baseDir, err := os.Getwd()
	assert(err)

	if !validateBridges(bridges) {
		os.Exit(1)
	}
	if !validateNumberOfVMs(createFlags.NumberOfVMs, bridges) {
		os.Exit(1)
	}

	if !createFlags.NoTmux {
		tmux.KillSession(createFlags.TMuxSessionName)
	}

	// create vms
	for i := 0; i < int(createFlags.NumberOfVMs); i++ {
		vm := createVM(i, hds, bridges, baseDir)

		// create disks
		for _, hd := range vm.hds {
			ensureHDExists(hd.path, createFlags.DiskSize)
		}

		qemuArgs := generateQEMUArgs(vm, createFlags.Memory, createFlags.Image)
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
