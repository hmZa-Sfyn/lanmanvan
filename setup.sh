#!/bin/bash

# Improved setup script for LanManVan
# Installs binary to ~/bin/lanmanvan and sets up aliases

set -e

# Create ~/bin if it doesn't exist and ensure it's in PATH
BIN_DIR="$HOME/bin"
mkdir -p "$BIN_DIR"

# Add ~/bin to PATH in .zshrc and .bashrc if not already there
for rc in "$HOME/.zshrc" "$HOME/.bashrc"; do
    if ! grep -q 'export PATH="$HOME/bin:$PATH"' "$rc" 2>/dev/null; then
        echo '' >> "$rc"
        echo '# Add ~/bin to PATH for local binaries' >> "$rc"
        echo 'export PATH="$HOME/bin:$PATH"' >> "$rc"
        echo "Added ~/bin to PATH in $rc"
    fi
done

echo "Running go mod tidy..."
go mod tidy

echo "Building lanmanvan..."
go build -o "$BIN_DIR/lanmanvan"

echo "Binary installed to $BIN_DIR/lanmanvan"

# Function to add or update alias
add_or_update_alias() {
    local rc_file=$1
    local alias_name=$2
    local target="lanmanvan"

    if ! grep -q "alias $alias_name=" "$rc_file" 2>/dev/null; then
        echo '' >> "$rc_file"
        echo "# LanManVan alias ($alias_name)" >> "$rc_file"
        echo "alias $alias_name='$target'" >> "$rc_file"
        echo "Added alias '$alias_name' in $rc_file"
    else
        sed -i "/alias $alias_name=/c\alias $alias_name='$target'" "$rc_file"
        echo "Updated alias '$alias_name' in $rc_file"
    fi
}

# Set up aliases in both shells
for rc in "$HOME/.zshrc" "$HOME/.bashrc"; do
    add_or_update_alias "$rc" "lanmanvan"
    add_or_update_alias "$rc" "lmvconsole"
done

echo ""
echo "Setup complete!"
echo ""
echo "The binary is now installed in ~/bin/lanmanvan and ~/bin is in your PATH."
echo "Both commands will work in any new terminal:"
echo "    lanmanvan"
echo "    lmvconsole"
echo ""
echo "To use them RIGHT NOW in this current shell, run:"
echo "    export PATH=\"\$HOME/bin:\$PATH\""
echo "    source ~/.zshrc   # (optional, reloads aliases too)"
echo ""
echo "After this, you can run your tool from anywhere with:"
echo "    lanmanvan"
echo "or"
echo "    lmvconsole"