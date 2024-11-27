#!/bin/bash

# Read the version from version.json
version=$(jq -r .version version.json)

# Update the Go file
awk -v version="$version" '{
    if ($0 ~ /helpers.AppConfig.VERSION/) {
        gsub(/helpers.AppConfig.VERSION/, "\"" version "\"");
    }
    print;
}' cmd/version/version.go > cmd/version/version_temp.go

# Replace the original file with the updated one
mv cmd/version/version_temp.go cmd/version/version.go
