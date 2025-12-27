#!/usr/bin/env python3
import os, sys
filepath = os.getenv('ARG_PATH', '')
content = os.getenv('ARG_CONTENT', '')
if not filepath:
    print('[!] Path required')
    sys.exit(1)

try:
    with open(filepath, 'w') as f:
        f.write(content)
    print(f'[+] File created: {filepath}')
    print(f'[+] Size: {len(content)} bytes')
except Exception as e:
    print(f'[!] Error: {e}')
