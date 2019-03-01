# GoPhish
Basic Experimental Golang phishing page.

This can be used on Android in an app called termux. Clone and build the repo,
and use Iptables to redirect all the trafic of your hotspot to the specified port.
Users will see whatever page you render, and might use their real credentials in
an attempt to gain internet access.

In termux, run android_setup.sh. Might not work on your mobile phone, but you might change the script as needed.

The idea is to create a mobile hotspot with the same SSID as a target open network. Most mobile phones will
only present users the strongest available wifi, and some mobile phones will disconnect from their current
wifi connection to connect to yours. After that, they will be presented what seems to be a captive portal
blocking the internet access until a successful login. You may choose a target user list, and if any of these
users logs in successfuly, the hotspot is turned off, allowing them to connect to the original open network with
normal internet access.
