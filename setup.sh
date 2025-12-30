#!/bin/bash
# Improved setup script for LanManVan CLI
# Installs binary, sets up aliases, copies modules folder
#
# Contributor: @l0n3ly_nat at x.com
# changes: optimized some things here and there. 15:32 - Tuesday
#
set -e
echo "🚀 Setting up LanManVan CLI..."
# Directories
BIN_DIR="$HOME/bin"
LANMANVAN_DIR="$HOME/lanmanvan"
MODULES_SRC="./modules" # Source: current dir has 'modules' folder
MODULES_DEST="$LANMANVAN_DIR/modules" # Destination: ~/lanmanvan/modules
# Create required directories
mkdir -p "$BIN_DIR"
mkdir -p "$LANMANVAN_DIR"
# Ensure ~/bin is in PATH
for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zprofile"; do
    if [ -f "$rc" ] && ! grep -q 'export PATH="$HOME/bin:$PATH"' "$rc"; then
        echo '' >> "$rc"
        echo '# LanManVan: Add ~/bin to PATH' >> "$rc"
        echo 'export PATH="$HOME/bin:$PATH"' >> "$rc"
        echo "Added ~/bin to PATH in $rc"
    fi
done
# Build the Go binary
echo "🔨 Building lanmanvan binary..."
go mod tidy
go build -o "$BIN_DIR/lanmanvan"
echo "✅ Binary installed to $BIN_DIR/lanmanvan"
# Copy modules directory if it exists
if [ -d "$MODULES_SRC" ]; then
    echo "📂 Copying modules directory to $MODULES_DEST..."
    rsync -av --delete "$MODULES_SRC/" "$MODULES_DEST/"
    echo "✅ Modules synced to $MODULES_DEST"
else
    echo "⚠️ Warning: './modules' directory not found. Skipping copy."
    echo " You can manually place modules in $MODULES_DEST later."
fi
# Function to safely add or update an alias (prevents duplicates)
add_or_update_alias() {
    local rc_file=$1
    local alias_name=$2
    local alias_command=$3  # Full command for the alias

    # Remove any existing alias with the same name (backup on macOS)
    if grep -q "alias $alias_name=" "$rc_file" 2>/dev/null; then
        sed -i.bak "/alias $alias_name=/d" "$rc_file" 2>/dev/null || \
        sed -i "" "/alias $alias_name=/d" "$rc_file"
        rm -f "${rc_file}.bak" 2>/dev/null
        echo "Updated existing alias '$alias_name' in $rc_file"
    fi

    # Add the new alias
    echo '' >> "$rc_file"
    echo "# LanManVan CLI alias" >> "$rc_file"
    echo "alias $alias_name='$alias_command'" >> "$rc_file"
    echo "Added alias '$alias_name' → '$alias_command' in $rc_file"
}
# Set up standard aliases
for rc in "$HOME/.zshrc" "$HOME/.bashrc" "$HOME/.bash_profile" "$HOME/.zprofile"; do
    if [ -f "$rc" ]; then
        add_or_update_alias "$rc" "lanmanvan" "lanmanvan -modules $MODULES_DEST"
        add_or_update_alias "$rc" "lmvconsole" "lanmanvan -modules $MODULES_DEST"
        # New alias: lmv_update – pulls latest repo and re-runs setup
        add_or_update_alias "$rc" "lmv_update" "cd /tmp && rm -rf lanmanvan && git clone https://github.com/hmZa-Sfyn/lanmanvan && cd lanmanvan && chmod +x setup.sh && ./setup.sh"
    fi
done
echo ""
echo "🎉 Setup complete!"
echo ""
echo "Modules location: $MODULES_DEST"
echo "Binary location: $BIN_DIR/lanmanvan"
echo ""
echo "You can now use these commands in any NEW terminal:"
echo " lanmanvan      # runs: lanmanvan -modules ~/lanmanvan/modules"
echo " lmvconsole    # same as above"
echo " lmv_update     # pulls latest version from GitHub and re-runs setup"
echo ""
echo "To use them IMMEDIATELY in this terminal, run:"
echo " export PATH=\"\$HOME/bin:\$PATH\""
echo " source ~/.zshrc 2>/dev/null || source ~/.bashrc 2>/dev/null || true"
echo ""
echo "After that, just type:"
echo " lanmanvan"
echo "Or to update in the future:"
echo " lmv_update"
echo "Enjoy your LanManVan CLI! 🚀"
