#!/bin/bash
CONFIG_FILE="/boot/firmware/config.txt"
BACKUP_FILE="/boot/firmware/config.txt.bak"

# Backup original file
sudo cp "$CONFIG_FILE" "$BACKUP_FILE"

# Define GPIO pins to set low
GPIO_PINS=(0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27)

# Check if GPIO lines already exist before adding them
for pin in "${GPIO_PINS[@]}"; do
    LINE="gpio=$pin=op,dl"
    
    # Check if the line already exists
    if ! grep -q "$LINE" "$CONFIG_FILE"; then
        echo "$LINE" | sudo tee -a "$CONFIG_FILE"
    fi
done

echo "GPIO pins set to low on startup (duplicates avoided)."
