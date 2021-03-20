# go-pfd

## Compilation
env GOOS=linux GOARCH=arm GOARM=6 go build .  
go-pfd  must be in /usr/local/bin/  

logs are in /var/log/bbox 
sudo mkdir /var/log/bbox  



## Raspberri Pi Zero W

### Influxdb

wget -qO- https://repos.influxdata.com/influxdb.key | sudo apt-key add -  
echo "deb https://repos.influxdata.com/debian buster stable" | sudo tee /etc/apt/sources.list.d/influxdb.list  
sudo apt update  
sudo apt install influxdb
sudo systemctl unmask influxdb  
sudo systemctl enable influxdb  

pi@bbox:~ $ dpkg-query -l 'influx*'  
Desired=Unknown/Install/Remove/Purge/Hold  
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend  
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)  
||/ Name           Version      Architecture Description  
+++-==============-============-============-=================================  
ii  influxdb       1.8.4-1      armhf        Distributed time-series database.  

go-pfd database default config:  
database: BBOX  
username: admin   
password: generic  

CREATE DATABASE BBOX  
CREATE USER admin WITH PASSWORD 'generic' WITH ALL PRIVILEGES  

### Grafana

wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -  
echo "deb https://packages.grafana.com/oss/deb stable main" | sudo tee -a /etc/apt/sources.list.d/grafana.list  
sudo apt update  
sudo apt-get install -y grafana  
sudo systemctl enable grafana-server  


pi@bbox:~ $ dpkg-query -l 'grafana-rpi'  
Desired=Unknown/Install/Remove/Purge/Hold  
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend  
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)  
||/ Name           Version      Architecture Description  
+++-==============-============-============-=================================  
ii  grafana-rpi    7.3.7        armhf        Grafana  


### Config bbox-small (/boot/config.txt)  
####CAN 12 MHz clock!  
dtoverlay=mcp2515-can0,oscillator=12000000,interrupt=25  
enable_uart=1  
start_x=0  

For permanent can0:
/etc/network/interfaces.d must hold the following can file  
auto can0  
iface can0 inet manual  
pre-up ip link set $IFACE type can bitrate 500000 listen-only off  
up /sbin/ifconfig $IFACE up  
down /sbin/ifconfig $IFACE down  

For debuging purpose optional following package can be installed:  
sudo apt-get install can-utils  


### Horloge GPS   
dtoverlay=disable-bt  
dtoverlay=pps-gpio,gpiopin=18  

### CMDLINE (/boot/cmdline.txt)  
root=PARTUUID=4cbd14f4-02 rootfstype=ext4 elevator=deadline fsck.repair=yes rootwait  

### RTC DS3231
Enable I2C  and reboot

sudo apt install python-smbus i2c-tools  
sudo i2cdetect -y 1 -> Should show id 68 present  

add /boot/config.txt and reboot  
dtoverlay=i2c-rtc,ds3231  

sudo apt -y remove fake-hwclock  
sudo update-rc.d -f fake-hwclock remove  

comment the 3 following lines in /lib/udev/hwclock-set:  

if [ -e /run/systemd/system ] ; then  
exit 0  
fi  

Check date time and if  not correct:  
sudo hwclock -w  

### systemd (/etc/systemd/system/bbox.service)  
[Unit]  
Description=Black Box Service
  
Wants=network.target  
After=syslog.target network-online.target  
  
[Service]  
Type=simple  
ExecStart=/usr/local/bin/go-pfd  
Restart=on-failure  
RestartSec=3  
KillMode=process  
  
[Install]  
WantedBy=multi-user.target  

sudo chmod 640 /etc/systemd/system/bbox.service  
sudo systemctl daemon-reload  
sudo systemctl enable bbox  

### RaspAP
curl -sL https://install.raspap.com | bash -s -- --yes --openvpn 0 --adblock 0  


Open the RaspAP admin interface in your browser, usually http://raspberrypi.local.  
The status widget should indicate that hostapd is inactive. This is expected.  
Confirm that the Wireless Client dashboard widget displays an active connection.  
Choose Hotspot > Advanced and enable the WiFi client AP mode option.  
Choose Save settings and Start hotspot.  
Wait a few moments and confirm that your AP has started.  

### DHCP bbox-small (/etc/dhcpcd.conf)
interface wlan0  
static ip_address=192.168.1.200/24  
static routers=192.168.1.1  
static domain_name_servers=192.168.1.1  