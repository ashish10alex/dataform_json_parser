#!/bin/bash

#!/bin/sh

# Exit immediately if a command exits with a non-zero status
set -e

# GitHub repository owner and name
REPO_OWNER="ashish10alex"
REPO_NAME="dj"

# Detect the operating system and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

echo $OS
echo $ARCH


# if os is darwin then change it to Darwin
# if os is linux then change it to Linux
case $OS in
    darwin)
        OS="Darwin"
        ;;
    linux)
        OS="Linux"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Translate architecture names to match GitHub release naming
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64 | arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Get the latest release download URL for the appropriate tar.gz file
RELEASE_URL=$(curl -s https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest | \
    grep "browser_download_url.*${OS}.*${ARCH}.*tar.gz" | cut -d '"' -f 4)


echo $RELEASE_URL

# Check if the URL is empty
if [ -z "$RELEASE_URL" ]; then
    echo "Could not find a release for OS: $OS and ARCH: $ARCH"
    exit 1
fi

# Temporary directory for download and extraction
TMP_DIR=$(mktemp -d)

# Change to the temporary directory
cd $TMP_DIR

# Download the latest release
curl -L -o release.tar.gz $RELEASE_URL

# Extract the tar.gz file
tar -xzvf release.tar.gz

# Find the binary (assuming it's in the root of the tar)
BINARY=$(find . -type f -perm +111 -exec basename {} \;)

# Check if the binary is found
if [ -z "$BINARY" ]; then
    echo "Could not find the binary file."
    exit 1
fi

# Move the binary to /usr/local/bin (requires sudo)
sudo mv $BINARY /usr/local/bin/

# Make sure the binary is executable
sudo chmod +x /usr/local/bin/$BINARY

# Clean up
cd -
rm -rf $TMP_DIR

echo "Installation completed. You can now use $BINARY."

