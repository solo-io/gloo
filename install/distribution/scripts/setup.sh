#!/bin/bash

echo "        _                  _   _ ";
echo "   __ _| | ___   ___   ___| |_| |";
echo "  / _\` | |/ _ \ / _ \ / __| __| |";
echo " | (_| | | (_) | (_) | (__| |_| |";
echo "  \__, |_|\___/ \___/ \___|\__|_|";
echo "  |___/                          ";
echo "";

# Get current script directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Get operating system name
if [[ "$(uname -s)" = "Darwin" ]]; then
  OS=darwin
else
  OS=linux
fi
echo "Detected operating system: $OS"

# Copy matching glooctl executable to user home directory
filename="glooctl-${OS}-amd64"
mkdir -p "$HOME/.gloo/bin"
cp "$DIR/$filename" "$HOME/.gloo/bin/glooctl"
chmod +x "$HOME/.gloo/bin/glooctl"
echo "Copying $DIR/$filename to $HOME/.gloo/bin/glooctl"

# Create symlink to executable
ln -sf "$HOME/.gloo/bin/glooctl" "/usr/local/bin/glooctl"
if [[ $? -eq 0 ]]; then
    echo "Creating symbolic link: /usr/local/bin/glooctl -> $HOME/.gloo/bin/glooctl"
    echo ""
    echo "glooctl should now be on your PATH! Try running 'glooctl --version' from your shell"
else
    echo "Could not create symbolic link in /usr/local/bin/glooctl"
    echo ""
    echo "Please add \"$HOME/.gloo/bin/glooctl\" to your PATH environment variable"
fi

