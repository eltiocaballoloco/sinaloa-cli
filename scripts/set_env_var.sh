#!/bin/bash

# Come back folder
cd ..

# Check if the .env file exists
if [ ! -f .env ]; then
    echo "The .env file does not exist."
    exit 1
fi

# Create a temporary file to store the exported variables
temp_env_file=$(mktemp)

# Read variables from .env file and export them to the temporary file
while IFS= read -r line; do
    if [[ ! -z "$line" ]]; then
        export "$line"
        echo "export $line" >> "$temp_env_file"
    fi
done < .env

# Source the temporary file to make the variables available globally
source "$temp_env_file"

# Clean up the temporary file
rm -f "$temp_env_file"

echo "[INFORMATION] Environment variables from '.env' file have been added to the system."
