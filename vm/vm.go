package vm

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type VMFlags struct {
	NumberOfVMs      uint8
	BridgeInterfaces string
	NoTmux           bool
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

var (
	UUIDs = []string{
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
)

func shiftSerial() (string, error) {
	if len(UUIDs) < 1 {
		return "", fmt.Errorf("No mac addresses left.")
	}
	serial := UUIDs[0]
	UUIDs = UUIDs[1:]
	return serial, nil
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

	var serial string
	serial, err = shiftSerial()
	if err != nil {
		return nil, err
	}

	vm := &VM{
		Serial:  serial,
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
