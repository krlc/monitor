#!/bin/bash

sudo useradd monitorservice -s /sbin/nologin -M
sudo usermod -a -G root monitorservice
sudo cp /home/monitor/monitorservice.service /lib/systemd/system/monitorservice.service
sudo chmod 755 /lib/systemd/system/monitorservice.service
sudo chmod -R 755 ./monitor_linux_x64
sudo chmod -R 777 ./visitors.log

sudo systemctl enable monitorservice.service
sudo systemctl start monitorservice
sudo journalctl -f -u monitorservice