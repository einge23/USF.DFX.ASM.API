#!/bin/bash

# Step 1: Create the service file
SERVICE_PATH="/etc/systemd/system/aacUi.service"
echo "[Unit]
Description=AAC UI Service
After=graphical.target

[Service]
ExecStart=/usr/bin/npm run dev
WorkingDirectory=/home/dfxp/Desktop/AutomatedAccessControl/Repos/USF.DFX.ASM.UI
Restart=always
User=dfxp
Environment=DISPLAY=:0
Environment=XAUTHORITY=/home/dfxp/.Xauthority

[Install]
WantedBy=graphical.target" | sudo tee $SERVICE_PATH > /dev/null

# Step 2: Reload systemd to recognize the new service
sudo systemctl daemon-reload

# Step 3: Enable the service to start on boot
sudo systemctl enable aacApi.service

# Step 4: Start the service manually
sudo systemctl start aacApi.service

# Step 5: Check the status of the service
sudo systemctl status aacApi.service
