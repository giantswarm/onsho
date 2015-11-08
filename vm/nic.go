package vm

import (
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
	MACs = []string{
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

func shiftMacAddress() (string, error) {
	if len(MACs) < 1 {
		return "", fmt.Errorf("No serial numbers left.")
	}
	mac := MACs[0]
	MACs = MACs[1:]
	return mac, nil
}

func bridgeExists(bridge string) bool {
	if _, err := os.Stat("/sys/class/net/" + bridge); os.IsNotExist(err) {
		return false
	}
	return true
}

func NewNIC(bridge string, bridgeOrderNumber int) (*NIC, error) {
	mac, err := shiftMacAddress()
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
