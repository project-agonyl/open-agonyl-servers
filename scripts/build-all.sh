#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to the project root directory
cd "$PROJECT_ROOT"

echo "Building all servers..."

# Iterate through all directories in the cmd folder
for server_dir in cmd/*/; do
    if [ -d "$server_dir" ]; then
        # Extract the server name from the directory path
        server_name=$(basename "$server_dir")
        echo "Building $server_name..."
        
        # Run the build script for this server
        ./scripts/build.sh "$server_name"
        
        if [ $? -eq 0 ]; then
            echo "$server_name build completed successfully!"
        else
            echo "ERROR: $server_name build failed!"
            exit 1
        fi
        echo ""
    fi
done

echo "All servers built successfully!"
