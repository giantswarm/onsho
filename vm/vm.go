package vm

import (
	"bytes"
	"crypto/rand"
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
	BridgeInterfaces string
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

type Population struct {
	Serials []string
	Macs    []string
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

func createMacAddress() (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("00:16:3e:%02x:%02x:%02x", buf[0], buf[1], buf[2]), nil
}

func NewVM(flags *VMFlags, baseDir string, pop *Population) (*VM, error) {
	imagePath, err := getImagePath(baseDir, flags.Image)
	if err != nil {
		return nil, err
	}

	vm := &VM{
		Memory:  flags.Memory,
		Image:   imagePath,
		BaseDir: baseDir,
	}

	if len(pop.Serials) > 0 {
		vm.Serial = pop.Serials[0]
		pop.Serials = pop.Serials[1:]
	} else {
		vm.Serial = uuid.NewV4().String()
	}
	hds := strings.Split(flags.HDs, ",")
	bridges := strings.Split(flags.BridgeInterfaces, ",")

	for _, device := range hds {
		hd, err := NewHD(device, flags.DiskSize, vm.BaseDir, vm.Serial)
		if err != nil {
			return vm, err
		}
		vm.HDs = append(vm.HDs, hd)
	}

	for i, bridge := range bridges {
		var mac string
		if len(pop.Macs) > 0 {
			mac = pop.Macs[0]
			pop.Macs = pop.Macs[1:]
		} else {
			mac, err = createMacAddress()
			if err != nil {
				return vm, err
			}
		}

		nic, err := NewNIC(bridge, i, mac)
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

func (vm *VM) Stop(tmuxSession string) error {
	return tmux.KillWindow(tmuxSession, vm.Serial)
}

func (vm *VM) Destroy() error {
	for _, hd := range vm.HDs {
		hd.Destroy()
	}

	return os.Remove(fmt.Sprintf("%s/machines/%s.json", vm.BaseDir, vm.Serial))
}

func StopAll(tmuxSession string) error {
	return tmux.KillSession(tmuxSession)
}

func DestroyAll(configDir string) error {
	err := os.RemoveAll(configDir + "/machines")
	if err != nil {
		return err
	}
	return os.RemoveAll(configDir + "/disks")
}

func (vm *VM) Wipe(flags *VMFlags) error {
	for _, hd := range vm.HDs {
		hd.Destroy()
	}

	vm.HDs = []*HD{}
	hds := strings.Split(flags.HDs, ",")

	for _, device := range hds {
		hd, err := NewHD(device, flags.DiskSize, vm.BaseDir, vm.Serial)
		if err != nil {
			return err
		}
		vm.HDs = append(vm.HDs, hd)
	}

	return nil
}

func (vm *VM) Status(tmuxSession string) (string, error) {
	windows, err := tmux.ListWindows(tmuxSession)
	if err != nil {
		return "error", err
	}

	status := "stopped"

	for _, w := range windows {
		if strings.Contains(w, vm.Serial) {
			status = "running"
		}
	}

	return status, nil
}
