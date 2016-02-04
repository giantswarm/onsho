# Moa

Manage QEMU VMs that will be provisioned by Mayu. With Moa you can replicate a Giant Swarm DC on a single machine or VM.

## Dependencies

 * `tmux`
 * `qemu` (>=2.4, including kvm support)
 * `linux`

## Setup

In order to give qemu cluster machines a separate network they can have fun in
you have to create a network bridge:
```
sudo brctl addbr bond0
sudo ip link set up dev bond0
sudo ip addr add 10.0.3.251/22 dev bond0
```

If you have systemd you can use systemd-networkd to create the bridge and make it remain after a reboot:
```
sudo cp host-conf/systemd-networkd/* /etc/systemd/network
sudo systemctl restart systemd-networkd
sudo systemctl enable systemd-networkd
```

To [allow qemu-bridge-helper](http://wiki.qemu.org/Features-Done/HelperNetworking#Setup) to manipulate our bridge add `allow bond0` to `/etc/qemu/bridge.conf`.

## Start a cluster

Start a qemu cluster:
```
$ ./moa create --num-vms=5 --image=ipxe/ipxe.iso # creates a 5 machine cluster
$ tmux a -t zoo                                # attach to the created tmux session
```

## Debugging

The simplest setup to get some output would be this:
```bash
$ ./moa create --num-vms=1 --image=ipxe/ipxe.iso --no-tmux
```
