package vm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type HD struct {
	Device string
	Path   string
	Size   string
}

func NewHD(device, size, baseDir, serial string) (*HD, error) {
	hd := &HD{
		Device: device,
		Path:   fmt.Sprintf("%s/disks/%s-%s.qcow2", baseDir, serial, device),
		Size:   size,
	}

	if _, err := os.Stat(hd.Path); os.IsNotExist(err) {
		cmd := exec.Command("qemu-img", "create", "-f", "qcow2", hd.Path, hd.Size)

		var stdout bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		return hd, cmd.Run()
	}

	return hd, nil
}
