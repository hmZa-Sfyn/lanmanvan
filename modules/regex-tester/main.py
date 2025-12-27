#!/usr/bin/env python3
import os, sys, re
pattern = os.getenv('ARG_PATTERN', '')
text = os.getenv('ARG_TEXT', '')
if not pattern or not text:
    print('[!] Pattern and text required')
    sys.exit(1)

try:
    matches = re.findall(pattern, text)
    print(f'[+] Pattern: {pattern}')
    print(f'[+] Matches found: {len(matches)}')
    for i, m in enumerate(matches[:10]):
        print(f'    {i+1}: {m}')
except Exception as e:
    print(f'[!] Error: {e}')
