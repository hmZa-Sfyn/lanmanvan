#!/bin/bash
# Create New Module Helper Script

if [ $# -lt 2 ]; then
    echo "Usage: ./create_module.sh <module_name> <type: python|bash>"
    echo ""
    echo "Example:"
    echo "  ./create_module.sh myrecon python"
    echo "  ./create_module.sh sysinfo bash"
    exit 1
fi

MODULE_NAME=$1
MODULE_TYPE=$2
MODULE_DIR="modules/$MODULE_NAME"

# Validate type
if [ "$MODULE_TYPE" != "python" ] && [ "$MODULE_TYPE" != "bash" ]; then
    echo "[!] Invalid type. Use 'python' or 'bash'"
    exit 1
fi

# Create directory
mkdir -p "$MODULE_DIR"
echo "[+] Created directory: $MODULE_DIR"

# Create module.yaml
cat > "$MODULE_DIR/module.yaml" << 'EOF'
name: $MODULE_NAME
description: "Description of your module"
type: $MODULE_TYPE
author: Your Name
version: 1.0.0
tags:
  - custom
options:
  target:
    type: string
    description: Target parameter
    required: true
required:
  - target
EOF

# Create main script based on type
if [ "$MODULE_TYPE" = "python" ]; then
    cat > "$MODULE_DIR/main.py" << 'EOF'
#!/usr/bin/env python3
"""
Module Description
"""

import os
import sys

def main():
    # Get arguments from environment variables
    target = os.getenv('ARG_TARGET') or 'localhost'
    
    print(f"[*] Module executing on {target}")
    
    try:
        # Your code here
        print("[+] Module completed successfully!")
    except Exception as e:
        print(f"[!] Error: {e}")
        sys.exit(1)

if __name__ == '__main__':
    main()
EOF
    chmod +x "$MODULE_DIR/main.py"
    echo "[+] Created main.py"
else
    cat > "$MODULE_DIR/main.sh" << 'EOF'
#!/bin/bash
# Module Description

TARGET="${ARG_TARGET:-localhost}"

echo "[*] Module executing on $TARGET"

# Your code here

echo "[+] Module completed successfully!"
EOF
    chmod +x "$MODULE_DIR/main.sh"
    echo "[+] Created main.sh"
fi

echo "[+] Module created successfully!"
echo ""
echo "Next steps:"
echo "1. Edit $MODULE_DIR/module.yaml to update metadata"
echo "2. Edit $MODULE_DIR/main.${MODULE_TYPE%%n*} to add your code"
echo "3. Run 'list' in LanManVan to see your new module"
