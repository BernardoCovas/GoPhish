#!/bin/sh

# This script is intended to be run on termux/android shell
# to setup iptables to redirect http/https requests to a custom
# server. This is perceived by clients as a captive portal.

su -c "sysctl net.ipv4.ip_forward=1 "
su -c "sysctl -w net.ipv4.conf.wlan0.route_localnet=1"
su -c "iptables \
    -t nat \
    -A PREROUTING \
    -p tcp \
    --match multiport --dports 80,443 \
    -j DNAT --to-destination 127.0.0.1:8080"