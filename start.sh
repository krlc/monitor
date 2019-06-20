#!/bin/bash

sudo systemctl start monitorservice
sudo journalctl -f -u monitorservice