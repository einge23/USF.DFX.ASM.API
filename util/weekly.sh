#!/bin/bash

# this file must be set in the raspberry pi builtin crontab
# sudo crontab -e to edit

#establish db
DB_FILE="/home/dfxp/Desktop/AutomatedAccessControl/Repos/USF.DFX.ASM.API/test.db"

# Query to fetch the last run date
FETCH_DATE_QUERY="SELECT last_ran_date FROM settings WHERE name = 'default';"
LAST_RAN_DATE=$(sqlite3 "$DB_FILE" "$FETCH_DATE_QUERY")

# Get the current date and current week number
CURRENT_DATE=$(date +%Y-%m-%d)
CURRENT_WEEK=$(date +%U)
LAST_RAN_WEEK=$(date -d "$LAST_RAN_DATE" +%U 2>/dev/null) # Use 'date -d' to parse the last run date

# Check if the script has already run this week
if [[ "$LAST_RAN_WEEK" == "$CURRENT_WEEK" ]]; then
    echo "Weekly hours have already been reset this week. Exiting."
    exit 0
fi

# Run the rest of the script
echo "Running weekly hour reset script for this week..."

FETCH_QUERY="SELECT default_user_weekly_hours FROM settings WHERE name = 'default';"

WEEKLY_HOURS=$(sqlite3 "$DB_FILE" "$FETCH_QUERY")

WEEKLY_MINUTES=$((WEEKLY_HOURS * 60))

UPDATE_QUERY="UPDATE users SET weekly_minutes = $WEEKLY_MINUTES;"

sqlite3 "$DB_FILE" "$UPDATE_QUERY"

# Update the last_ran_date to today
UPDATE_QUERY="UPDATE settings SET last_ran_date = '$CURRENT_DATE' WHERE name = 'default';"
sqlite3 "$DB_FILE" "$UPDATE_QUERY"

echo "Weekly hours reset and last_ran_date updated."
