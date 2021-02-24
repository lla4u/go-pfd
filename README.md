# go-pfd

## Compilation
env GOOS=linux GOARCH=arm GOARM=6 go build .  

## RaspBerri Zero

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
