package vm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/giantswarm/moa/tmux"
	"github.com/satori/go.uuid"
)

type VMFlags struct {
	NumberOfVMs      uint8
	BridgeInterfaces string
	NoTMux           bool
	TMuxSessionName  string
	HDs              string
	Memory           uint16
	DiskSize         string
	Image            string
}

type VM struct {
	HDs     []*HD
	NICs    []*NIC
	Serial  string
	Memory  uint16
	Image   string
	BaseDir string
}

func getImagePath(baseDir, image string) (string, error) {
	var imagePath string
	if strings.Contains(image, "/") {
		imagePath = image
	} else {
		imagePath = baseDir + "/images/" + image
	}
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return imagePath, fmt.Errorf("Image not found: %s", imagePath)
	}

	return imagePath, nil
}

func NewVM(flags *VMFlags, baseDir string) (*VM, error) {
	imagePath, err := getImagePath(baseDir, flags.Image)
	if err != nil {
		return nil, err
	}

	vm := &VM{
		Serial:  uuid.NewV4().String(),
		Memory:  flags.Memory,
		Image:   imagePath,
		BaseDir: baseDir,
	}

	hds := strings.Split(flags.HDs, ",")
	bridges := strings.Split(flags.BridgeInterfaces, ",")

	for _, device := range hds {
		hd, err := NewHD(device, flags.DiskSize, baseDir, vm.Serial)
		if err != nil {
			return vm, err
		}
		vm.HDs = append(vm.HDs, hd)
	}

	for i, bridge := range bridges {
		nic, err := NewNIC(bridge, i)
		if err != nil {
			return vm, err
		}
		vm.NICs = append(vm.NICs, nic)
	}

	if err := vm.Save(); err != nil {
		return vm, err
	}

	return vm, nil
}

func (vm *VM) Save() error {
	err := saveJson(vm, fmt.Sprintf("%s/machines/%s.json", vm.BaseDir, vm.Serial))
	if err != nil {
		return err
	}
	return nil
}

func Load(baseDir, serial string) (*VM, error) {
	vm := &VM{}
	err := loadJson(vm, fmt.Sprintf("%s/machines/%s.json", baseDir, serial))
	if err != nil {
		return nil, err
	}

	return vm, nil
}

func List(baseDir string) ([]*VM, error) {
	files, err := ioutil.ReadDir(baseDir + "/machines")
	if err != nil {
		return nil, err
	}

	vms := []*VM{}

	for _, f := range files {
		vm, err := Load(baseDir, strings.Replace(f.Name(), ".json", "", 1))
		if err != nil {
			return vms, err
		}

		vms = append(vms, vm)
	}

	return vms, nil
}

func (vm *VM) GenerateQEMUArgs() []string {
	var args []string
	args = append(args, "-m", strconv.Itoa(int(vm.Memory)))
	args = append(args, "-cdrom", vm.Image)
	args = append(args, "-boot", "d")
	args = append(args, "-nographic")
	args = append(args, "-device", "sga")
	args = append(args, "-serial", "mon:stdio")

	for vlan, nic := range vm.NICs {
		args = append(args, "-net", fmt.Sprintf("bridge,br=%s,vlan=%d", nic.Bridge, vlan))
		args = append(args, "-net", fmt.Sprintf("nic,vlan=%d,model=%s,macaddr=%s,addr=%s", vlan, nic.Dev, nic.Mac, nic.Addr))
	}

	args = append(args, "--cpu", "host")
	args = append(args, "-enable-kvm")

	for _, hd := range vm.HDs {
		args = append(args, fmt.Sprintf("-%s", hd.Device), hd.Path)
	}

	args = append(args, "-smbios", fmt.Sprintf("type=3,serial=%s", vm.Serial))
	args = append(args, "-monitor", fmt.Sprintf("unix:/tmp/tinyswarm-%s.sock,server,nowait", vm.Serial))

	return args
}

func (vm *VM) Start(tmuxSession string, noTMux bool) {
	qemuArgs := vm.GenerateQEMUArgs()
	qemuCmd := fmt.Sprintf("%s %s", "qemu-system-x86_64", strings.Join(qemuArgs, " "))

	if !noTMux {
		tmux.NewWindow(tmuxSession, vm.Serial, qemuCmd)
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

func (vm *VM) Stop() {
	tmux.KillWindow(vm.Serial)
}

func (vm *VM) Destroy() {
	for _, hd := range vm.HDs {
		os.Remove(hd.Path)
	}
	os.Remove(fmt.Sprintf("%s/machines/%s.json", vm.BaseDir, vm.Serial))
}

func StopAll(tmuxSession string) {
	tmux.KillSession(tmuxSession)
}

func DestroyAll(configDir string) {
	os.RemoveAll(configDir + "/machines")
	os.RemoveAll(configDir + "/disks")
}
