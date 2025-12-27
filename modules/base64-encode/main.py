#!/usr/bin/env python3
import base64
import os
import sys

data = os.getenv('ARG_DATA', '')
if not data:
    print('[!] Error: data is required')
    sys.exit(1)

try:
    encoded = base64.b64encode(data.encode()).decode()
    print(f'[+] Original: {data}')
    print(f'[+] Encoded:  {encoded}')
except Exception as e:
    print(f'[!] Error: {e}')
    sys.exit(1)
