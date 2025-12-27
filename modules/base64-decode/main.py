#!/usr/bin/env python3
import base64
import os
import sys

data = os.getenv('ARG_DATA', '')
if not data:
    print('[!] Error: data is required')
    sys.exit(1)

try:
    decoded = base64.b64decode(data).decode()
    print(f'[+] Encoded: {data}')
    print(f'[+] Decoded: {decoded}')
except Exception as e:
    print(f'[!] Error: Invalid Base64 data - {e}')
    sys.exit(1)
