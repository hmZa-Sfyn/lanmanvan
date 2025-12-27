#!/usr/bin/env python3
import urllib.parse
import os
import sys

data = os.getenv('ARG_DATA', '')
if not data:
    print('[!] Error: data is required')
    sys.exit(1)

try:
    decoded = urllib.parse.unquote(data)
    print(f'[+] Encoded: {data}')
    print(f'[+] Decoded: {decoded}')
except Exception as e:
    print(f'[!] Error: {e}')
    sys.exit(1)
