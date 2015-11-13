package vm

import (
	"crypto/rand"
	"fmt"
	"os"
)

type NIC struct {
	Bridge string
	Mac    string
	Dev    string
	Addr   string
}

var (
	networkInterfaceAddresses = []string{
		"03", "05", "07",
	}
)

func getMacAddress() (string, error) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("00:16:3e:%02x:%02x:%02x", buf[0], buf[1], buf[2]), nil
}

func bridgeExists(bridge string) bool {
	if _, err := os.Stat("/sys/class/net/" + bridge); os.IsNotExist(err) {
		return false
	}
	return true
}

func NewNIC(bridge string, bridgeOrderNumber int) (*NIC, error) {
	mac, err := getMacAddress()
	if err != nil {
		return nil, err
	}

	if !bridgeExists(bridge) {
		return nil, fmt.Errorf(`Interface %s doesn't exist, you might want to set it up:

sysctl net.ipv4.ip_forward=1

brctl addbr %s
ip add add cidr_addr dev %s
ip link set %s up
`, bridge, bridge, bridge, bridge)
	}

	return &NIC{
		Bridge: bridge,
		Mac:    mac,
		Dev:    "e1000",
		Addr:   networkInterfaceAddresses[bridgeOrderNumber],
	}, nil
}
