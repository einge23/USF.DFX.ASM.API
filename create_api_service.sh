#!/bin/bash

# Step 1: Create the service file
SERVICE_PATH="/etc/systemd/system/aacApi.service"
echo "[Unit]
Description=AAC Go API Service
After=network.target

[Service]
ExecStart=/home/dfxp/Desktop/AutomatedAccessControl/aacApiRun.sh
WorkingDirectory=/home/dfxp/Desktop/AutomatedAccessControl/Repos/USF.DFX.ASM.API
Restart=always
User=dfxp
Environment=DISPLAY=:0
Environment=XAUTHORITY=/home/dfxp/.Xauthority

[Install]
WantedBy=multi-user.target" | sudo tee $SERVICE_PATH > /dev/null

# Step 2: Reload systemd to recognize the new service
sudo systemctl daemon-reload

# Step 3: Enable the service to start on boot
sudo systemctl enable aacApi.service

# Step 4: Start the service manually
sudo systemctl start aacApi.service

# Step 5: Check the status of the service
sudo systemctl status aacApi.service
