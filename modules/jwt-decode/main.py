#!/usr/bin/env python3
import os, sys, json, base64
token = os.getenv('ARG_TOKEN', '')
parts = token.split('.')
if len(parts) != 3:
    print('[!] Invalid JWT')
    sys.exit(1)
try:
    header = json.loads(base64.urlsafe_b64decode(parts[0] + '=='))
    payload = json.loads(base64.urlsafe_b64decode(parts[1] + '=='))
    print('[+] Header:', json.dumps(header, indent=2))
    print('[+] Payload:', json.dumps(payload, indent=2))
except Exception as e:
    print(f'[!] Error: {e}')
