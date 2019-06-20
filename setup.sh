#!/bin/bash

sudo useradd monitorservice -s /sbin/nologin -M
sudo usermod -a -G root monitorservice
sudo cp /home/monitor/monitorservice.service /lib/systemd/system/monitorservice.service
sudo chmod 755 /lib/systemd/system/monitorservice.service

sudo systemctl enable monitorservice.service
sudo systemctl start monitorservice
sudo journalctl -f -u monitorservice