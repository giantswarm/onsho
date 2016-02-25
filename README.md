# Onsho

[![](https://godoc.org/github.com/giantswarm/onsho?status.svg)](http://godoc.org/github.com/giantswarm/onsho) [![IRC Channel](https://img.shields.io/badge/irc-%23giantswarm-blue.svg)](https://kiwiirc.com/client/irc.freenode.net/#giantswarm)

Manage QEMU VMs that will be provisioned by [Mayu](https://github.com/giantswarm/mayu). With Onsho you can replicate a Giant Swarm DC on a single machine or VM.

## Prerequisites

## Getting Onsho

Download the latest release from here: https://github.com/giantswarm/onsho/releases/latest

Clone the latest git repository version from here: git@github.com:giantswarm/onsho.git

### How to build

#### Dependencies

 * `tmux`
 * `qemu` (>=2.4, including kvm support)
 * `linux`

#### Building the standard way

```
make && sudo make install
```

## Running Onsho

### Setup

In order to give qemu cluster machines a separate network they can have fun in
you have to create a network bridge:

```
sudo brctl addbr onsho0
sudo ip link set up dev onsho0
sudo ip addr add 10.0.3.251/22 dev onsho0
```

If you have systemd you can use systemd-networkd to create the bridge and make it remain after a reboot:

```
sudo cp host-conf/systemd-networkd/* /etc/systemd/network
sudo systemctl restart systemd-networkd
sudo systemctl enable systemd-networkd
```

To [allow qemu-bridge-helper](http://wiki.qemu.org/Features-Done/HelperNetworking#Setup) to manipulate our bridge add `allow onsho0` to `/etc/qemu/bridge.conf`.

### Start a cluster

Start a qemu cluster:
```
$ onsho create --num-vms=5 --image=ipxe/ipxe.iso # creates a 5 machine cluster
$ tmux a -t zoo                                # attach to the created tmux session
```

## Debugging

The simplest setup to get some output would be this:
```bash
$ onsho create --num-vms=1 --image=ipxe/ipxe.iso --no-tmux
```

## Further Steps

Check more detailed documentation: [docs](docs)

Check code documentation: [godoc](https://godoc.org/github.com/giantswarm/onsho)

## Future Development

- Adapt Onsho to also work with VirtualBox on MacOS.

## Contact

- Mailing list: [giantswarm](https://groups.google.com/forum/!forum/giantswarm)
- IRC: #[giantswarm](irc://irc.freenode.org:6667/#giantswarm) on freenode.org
- Bugs: [issues](https://github.com/giantswarm/onsho/issues)

## Contributing & Reporting Bugs

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches, the contribution workflow as well as reporting bugs.

## License

Onsho is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.

## Origin of the Name

`onsho` (おんしょう[温床] pronounced "on-shou") is Japanese for breeding ground.
