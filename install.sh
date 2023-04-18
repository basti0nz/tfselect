#!/bin/bash

BIN_PATH='/usr/local/bin'

if ! command -v curl &> /dev/null; then
    echo "curl is not installed, please install it and try again."
    exit 1
fi

OS=`/usr/bin/uname -s`
ARCH=$(uname -m)
if [[ $OS == "Darwin" ]]; then
    PLATFORM="darwin_all"
elif [[ $OS == "Linux" ]]; then
  if [[ $ARCH == "x86_64" ]]; then
    PLATFORM="linux"
  elif [[ $ARCH == "aarch64" ]]; then
    PLATFORM="linux_arm64"
  else
    echo "Unsupported architecture: $ARCH"
    exit 1
  fi
else
  echo "Unsupported operating system: $OS"
  exit 1
fi

LATEST_RELEASE=`/usr/bin/curl -s https://api.github.com/repos/basti0nz/tfselect/releases/latest \
  | /usr/bin/grep "browser_download_url.*$PLATFORM.*amd64.tar.gz" \
  | /usr/bin/cut -d '"' -f 4`


echo "Downloading $LATEST_RELEASE"
curl -L -o tfselect.tar.gz $LATEST_RELEASE
tar -xzf tfselect.tar.gz
rm tfselect.tar.gz
mv tfselect $BIN_PATH/tfselect
chmod 0755 $BIN_PATH/tfselect

echo "tfselect installed successfully!"
