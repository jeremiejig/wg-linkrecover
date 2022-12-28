# wg-linkrecover

The goal of this program is to monitor the wireguard interface and detect any firewall session drop.\
To recover from a firewall session drop we just change the listen-port of the wireguard tunnel.

## Requirement

This code target only Linux kernel module for now.

Wireguard-monitor requires *root* or `cap_net_admin` to work.

## Installation

To install the program you can follow this step:

* `go get github.com/jeremiejig/wg-linkrecover/cmd/wg-linkrecover`
* optionnaly do `setcap cap_net_admin=ep wg-linkrecover`

## Launch

To start monitoring a wireguard link start the command as follow:

* `wg-linkrecover -ifname wg0`

Or create a systemd service:

```systemd
# /etc/systemd/system/wg-linkrecover@.service
[Unit]
Description=Wireguard link recover for %I
After=wg-quick@%i.service
Wants=wg-quick@%i.service
PartOf=wg-quick@%i.service
# Or
After=sys-devices-virtual-net-%i.device
Wants=sys-devices-virtual-net-%i.device
PartOf=sys-devices-virtual-net-%i.device

[Service]
RootDirectory=/usr/local
ExecStart=/bin/wg-linkrecover -ifname %I
DynamicUser=yes
ProtectProc=invisible
ProcSubset=pid
UMask=0077

AmbientCapabilities=cap_net_admin
CapabilityBoundingSet=cap_net_admin
SecureBits=no-setuid-fixup-locked

RestrictNamespaces=yes
ProtectHome=yes
ProtectHostname=yes
ProtectKernelTunables=yes
ProtectClock=yes
ProtectKernelLogs=yes
ProtectControlGroups=yes
ProtectKernelModules=yes
MemoryDenyWriteExecute=yes
RestrictRealtime=yes

PrivateDevices=yes
PrivateNetwork=no
IPAddressDeny=any
RestrictAddressFamilies=AF_NETLINK
PrivateIPC=yes
LockPersonality=yes
# Go specific
SystemCallFilter=@default @signal @io-event @basic-io clone openat pipe2 fcntl
# app specific
SystemCallFilter=@network-io
SystemCallErrorNumber=EPERM
SystemCallArchitectures=native

[Install]
WantedBy=multi-user.target
```
