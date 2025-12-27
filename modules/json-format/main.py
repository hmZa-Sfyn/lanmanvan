#!/usr/bin/env python3
import os, sys, json
data = os.getenv('ARG_DATA', '')
indent = int(os.getenv('ARG_INDENT', '2'))
try:
    parsed = json.loads(data)
    formatted = json.dumps(parsed, indent=indent)
    print('[+] Valid JSON')
    print(formatted)
except json.JSONDecodeError as e:
    print(f'[!] Invalid JSON: {e}')
