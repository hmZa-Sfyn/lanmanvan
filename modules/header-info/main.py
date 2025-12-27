#!/usr/bin/env python3
import os, sys, requests
url = os.getenv('ARG_URL', '')
if not url:
    print('[!] URL required')
    sys.exit(1)

try:
    r = requests.head(url, allow_redirects=True, timeout=5)
    print(f'[+] Status: {r.status_code}')
    print('[+] Headers:')
    for k, v in r.headers.items():
        print(f'    {k}: {v}')
except Exception as e:
    print(f'[!] Error: {e}')
