#!/bin/bash

# Script to tidy, build LanManVan, and set up 'lanmanvan' alias

set -e  # Exit on any error

PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY_PATH="$HOME/.lanmanvan"

echo "Running go mod tidy..."
go mod tidy

echo "Building lanmanvan..."
go build -o "$BINARY_PATH"

echo "Binary built and placed at $BINARY_PATH"

# Add/update alias in .zshrc
ZSHRC="$HOME/.zshrc"
if ! grep -q "alias lanmanvan=" "$ZSHRC" 2>/dev/null; then
    echo "" >> "$ZSHRC"
    echo "# LanManVan console alias" >> "$ZSHRC"
    echo "alias lanmanvan='$BINARY_PATH'" >> "$ZSHRC"
    echo "Added 'lanmanvan' alias to $ZSHRC"
else
    sed -i "/alias lanmanvan=/c\alias lanmanvan='$BINARY_PATH'" "$ZSHRC"
    echo "Updated existing 'lanmanvan' alias in $ZSHRC"
fi

# Add/update alias in .bashrc
BASHRC="$HOME/.bashrc"
if ! grep -q "alias lanmanvan=" "$BASHRC" 2>/dev/null; then
    echo "" >> "$BASHRC"
    echo "# LanManVan console alias" >> "$BASHRC"
    echo "alias lanmanvan='$BINARY_PATH'" >> "$BASHRC"
    echo "Added 'lanmanvan' alias to $BASHRC"
else
    sed -i "/alias lanmanvan=/c\alias lanmanvan='$BINARY_PATH'" "$BASHRC"
    echo "Updated existing 'lanmanvan' alias in $BASHRC"
fi

echo ""
echo "Setup complete!"
echo "To use the alias immediately in the current shell:"
echo "    source ~/.zshrc   # or source ~/.bashrc if you're using bash"
echo ""
echo "Now you can run the console from anywhere with:"
echo "    lanmanvan"
