curr_dir=$(pwd)

cd /
sysctl net.ipv4.ip_forward=1 
iptables \
    -t nat \
    -A PREROUTING \
    -p tcp \
    --match multiport --dports 80,443 \
    -j DNAT --to-destination 127.0.0.1:8080

sysctl -w net.ipv4.conf.wlan0.route_localnet=1
cd $curr_dir
