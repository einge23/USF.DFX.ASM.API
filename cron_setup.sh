#!/bin/bash

# Define the cron jobs you want
CRON_JOBS=$(cat <<EOF
# Run every monday morning at 8:00 AM
0 8 * * 1 echo "---------------- \$(date '+\\%Y-\\%m-\\%d \\%H:\\%M:\\%S') ----------------" >> /home/dfxp/Desktop/AutomatedAccessControl/cronLogs/weeklyLog.txt && /home/dfxp/Desktop/AutomatedAccessControl/Repos/USF.DFX.ASM.API/util/weekly.sh >> /home/dfxp/Desktop/AutomatedAccessControl/cronLogs/weeklyLog.txt 2>&1

# Run on reboot
@reboot echo "---------------- \$(date '+\\%Y-\\%m-\\%d \\%H:\\%M:\\%S') ----------------" >> /home/dfxp/Desktop/AutomatedAccessControl/cronLogs/weeklyLog.txt && /home/dfxp/Desktop/AutomatedAccessControl/Repos/USF.DFX.ASM.API/util/weekly.sh >> /home/dfxp/Desktop/AutomatedAccessControl/cronLogs/weeklyLog.txt 2>&1
EOF
)

# Backup the current system-wide crontab
sudo crontab -l > /tmp/current_crontab 2>/dev/null

# Combine the current crontab with the new cron jobs, avoiding duplicates
(echo "$CRON_JOBS"; cat /tmp/current_crontab) | sort | uniq | sudo crontab -

# Cleanup the temporary file
rm -f /tmp/current_crontab

echo "System-wide crontab updated successfully!"
