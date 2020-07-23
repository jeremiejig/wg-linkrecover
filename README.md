# Wireguard-Monitor

The goal of this program is to monitor the wireguard interface and detect any firewall session drop.\
To recover from a firewall session drop we just change the listen-port of the wireguard tunnel.

## Requirement

This code target only Linux kernel module for now.

Wireguard-monitor requires *root* or `cap_net_admin` + `cap_net_raw to work`.