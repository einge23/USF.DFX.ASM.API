# Script detects any mounted usb partitions, and removes them all.

#!/bin/bash

# Function to detect all mounted removable USB partitions
detect_mounted_usb_partitions() {
    # Use lsblk to find mounted removable partitions
    PARTITIONS=$(lsblk -ln -o NAME,TYPE,RM,MOUNTPOINT | grep "part" | awk '$2 == "part" && $3 == "1" && $4 != "" {print "/dev/" $1}')
    echo "$PARTITIONS"
}

# Function to unmount removable USB partitions
unmount_usb_partitions() {
    PARTITIONS=$(detect_mounted_usb_partitions)

    if [ -z "$PARTITIONS" ]; then
        echo "No mounted removable USB partitions detected."
        exit 0
    fi

    echo "Detected mounted removable partitions:"
    echo "$PARTITIONS"

    for PARTITION in $PARTITIONS; do
        echo "Unmounting $PARTITION..."
        umount "$PARTITION"
        if [ $? -eq 0 ]; then
            echo "$PARTITION unmounted successfully."
        else
            echo "Failed to unmount $PARTITION."
        fi
    done
}

# Execute the unmounting function
unmount_usb_partitions
