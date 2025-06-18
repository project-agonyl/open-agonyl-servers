#!/bin/bash

# Default binary name
binary_name="login-server"

# Check if a binary name is provided as an argument
if [ $# -eq 1 ]; then
    binary_name=$1
fi

# Set the Go environment variables for building for Windows
export GOARCH=amd64
export GOOS=windows

# Build for Windows
echo "Building $binary_name for Windows..."
go build -ldflags="-w -s" -o "bin/${binary_name}.exe" "cmd/${binary_name}/main.go"

# Reset Go environment variables to their defaults
export GOARCH=
export GOOS=

# Build for Linux
echo "Building $binary_name for Linux..."
go build -ldflags="-w -s" -o "bin/${binary_name}" "cmd/${binary_name}/main.go"

echo "$binary_name build complete!"
