#!/bin/bash

# Function to detect all mounted removable USB partitions
detect_mounted_usb_partitions() {
    # Use lsblk to find mounted removable partitions and exclude the SD card (/dev/mmcblk0*)
    PARTITIONS=$(lsblk -ln -o NAME,TYPE,RM,MOUNTPOINT | grep "part" | awk '$2 == "part" && $3 == "1" && $4 != "/" {print "/dev/" $1}')
    echo "$PARTITIONS"
}

# Function to unmount and eject removable USB partitions
unmount_and_eject_usb_partitions() {
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
            echo "Ejecting $PARTITION..."
            eject "$PARTITION"
            if [ $? -eq 0 ]; then
                echo "$PARTITION ejected successfully."
            else
                echo "Failed to eject $PARTITION."
            fi
        else
            echo "Failed to unmount $PARTITION."
        fi
    done
}

# Execute the unmount and eject function
unmount_and_eject_usb_partitions