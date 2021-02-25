# go-pfd

## Compilation
env GOOS=linux GOARCH=arm GOARM=6 go build .  

## Raspberri Pi Zero W

### Influxdb

pi@bbox:~ $ dpkg-query -l 'influx*'  
Desired=Unknown/Install/Remove/Purge/Hold  
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend  
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)  
||/ Name           Version      Architecture Description  
+++-==============-============-============-=================================  
ii  influxdb       1.8.4-1      armhf        Distributed time-series database.  

default config:  
database: bbox  
username: admin  
password: generic

CREATE DATABASE bbox  
CREATE USER admin WITH PASSWORD 'generic' WITH ALL PRIVILEGES  

Install
### Grafana
pi@bbox:~ $ dpkg-query -l 'grafana-rpi'  
Desired=Unknown/Install/Remove/Purge/Hold  
| Status=Not/Inst/Conf-files/Unpacked/halF-conf/Half-inst/trig-aWait/Trig-pend  
|/ Err?=(none)/Reinst-required (Status,Err: uppercase=bad)  
||/ Name           Version      Architecture Description  
+++-==============-============-============-=================================  
ii  grafana-rpi    7.3.7        armhf        Grafana  


### Config (/boot/config.txt)  
####CAN 12 MHz clock!  
dtoverlay=mcp2515-can0,oscillator=12000000,interrupt=25  
enable_uart=1  
start_x=0  

### Horloge GPS   
dtoverlay=disable-bt  
dtoverlay=pps-gpio,gpiopin=18  

### CMDLINE (/boot/cmdline.txt)  
root=PARTUUID=4cbd14f4-02 rootfstype=ext4 elevator=deadline fsck.repair=yes rootwait  

### DHCP (/etc/dhcpcd.conf)
interface wlan0  
static ip_address=192.168.1.200/24  
static routers=192.168.1.1  
static domain_name_servers=192.168.1.1  

### systemd (/etc/systemd/system/bbox.service)  
[Unit]  
Description=BBOX 
  
Wants=network.target  
After=syslog.target network-online.target  
  
[Service]  
Type=simple  
ExecStart=/usr/local/bin/go-pfd  
Restart=on-failure  
RestartSec=10  
KillMode=process  
  
[Install]  
WantedBy=multi-user.target  

sudo chmod 640 /etc/systemd/system/bbox.service  
sudo systemctl daemon-reload  
sudo systemctl enable bbox