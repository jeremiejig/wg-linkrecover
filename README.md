# wg-linkrecover

The goal of this program is to monitor the wireguard interface and detect any firewall session drop.\
To recover from a firewall session drop we just change the listen-port of the wireguard tunnel.

## Requirement

This code target only Linux kernel module for now.

Wireguard-monitor requires *root* or `cap_net_admin` to work.

## Installation

To install the program you can follow this step:

* `go get github.com/jeremiejig/wg-linkrecover`
* optionnaly do `setcap cap_net_admin=ep wg-linkrecover`

## Launch

To start monitoring a wireguard link start the command as follow:

* `wg-linkrecover -ifname wg0`

Or create a systemd service:

```systemd
[Unit]
Description=Wireguard link recover for wg0
After=wg-quick@wg0.service
Wants=wg-quick@wg0.service
PartOf=wg-quick@wg0.service

[Service]
ExecStart=/usr/local/bin/wg-linkrecover
AmbientCapabilities=cap_net_admin

[Install]
WantedBy=multi-user.target
```
