# Getting Started with Onsho

Onsho helps you to simulate a bare metal iPXE/PXE environment on a single test machine. This could be your laptop or some other physical machine that is available to you.

For now Onsho only supports QEMU/KVM. Which makes Linux a requirement. But you should be able to use VMWare Fusion on a Mac too. The main challenge is that Onsho needs nested virtualization that currently isn't possible with VirtualBox. We have also created similar setups with just booting VirtualBox instances but this is currently not integrated into Onsho.

## Prerequisites

 * You need a Linux box
 * Install QEMU incl. KVM support (eg. pacman -Sy qemu, apt-get install qemu-kvm)
 * Install tmux
 * Install docker (if you want to use Mayu to create a CoreOS cluster)

## Building Onsho

```
git clone https://github.com/giantswarm/onsho.git
cd onsho
make && sudo make install
```

## Start mayu

Fetch a Mayu release and create your own configuration and docker image:

```
wget https://downloads.giantswarm.io/mayu/latest/mayu.tar.gz
mkdir mayu
tar xzf mayu.tar.gz -C mayu
cd mayu
```

Fetch the CoreOS version you would like to use:

```
./fetch-coreos-image 835.13.0
```

Check the versions of docker, etcd and fleet you would like to install. There are defaults defined in the ´./fetch-mayu-asset` script.

```
grep 'VERSION=' fetch-mayu-assets
./fetch-mayu-assets
```

There are a few things you should configure before Mayu is usable:

 * add your SSH key to the config.yaml (replace '<your public key>')
 * adapt docker, fleet and etcd versions
 * add `no_secure: true` to the config if you don't want to use TLS in your local setup
 * change the interface (bond0) to something like onsho0

```
cp config.yaml.dist config.yaml
vi config.yaml
```

Now Mayus configuration is in place and we can build and run the container.

```
docker build -t mayu .

mkdir cluster
docker run --cap-add=NET_ADMIN --net=host -v $(pwd)/cluster:/var/lib/mayu -v $(pwd)/images:/usr/lib/mayu/images -v $(pwd)/assets:/usr/lib/mayu/assets --name=mayu mayu:latest
```

## Configure Onsho

In order to give qemu cluster machines a separate network they can have fun in you have to create a network bridge:

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

## Start a cluster

Now with Mayu and the bridge configured you can start your local CoreOS cluster.

```
onsho create --num-vms=3
tmux a -t zoo
```

## Debugging

The simplest setup to get some output would be this:
```bash
$ onsho create --num-vms=1 --image=ipxe/ipxe.iso --no-tmux
```
