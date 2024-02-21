#!/bin/bash

# Define the current version
htmxVersion="1.9.10"

# Define the output directory
outputDir="static"

# Define the URL
url="https://unpkg.com/htmx.org@$htmxVersion/dist/htmx.min.js"

# Download the file with progress
echo "Downloading htmx..."
curl -L $url -o $outputDir/htmx.min.js --progress-bar

# Check if the download was successful
if [ $? -ne 0 ]; then
    echo "Download failed."
    exit 1
fi

# Fetch the latest version
echo "Fetching the latest version..."
latestVersion=$(curl -s https://unpkg.com/htmx.org | grep -oP 'htmx.org@\K[^y]*')

# Check if fetching the latest version was successful
if [ $? -ne 0 ]; then
    echo "Failed to fetch the latest version."
    exit 1
fi

# Display the latest version
echo "Latest version of htmx is: $latestVersion"

# Check if a new version is available
if [ "$htmxVersion" != "$latestVersion" ]; then
    echo "A new version of htmx is available."
    # Add code here to download and update to the new version
fi