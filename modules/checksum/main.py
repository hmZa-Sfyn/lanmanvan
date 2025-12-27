#!/usr/bin/env python3
import os, sys, hashlib
filepath = os.getenv('ARG_FILE', '')
hashtype = os.getenv('ARG_TYPE', 'sha256').lower()

if not filepath or not os.path.exists(filepath):
    print('[!] File not found')
    sys.exit(1)

if hashtype not in ['md5', 'sha1', 'sha256']:
    print('[!] Invalid hash type')
    sys.exit(1)

hash_obj = hashlib.new(hashtype)
with open(filepath, 'rb') as f:
    for chunk in iter(lambda: f.read(4096), b''):
        hash_obj.update(chunk)

print(f'[+] File: {filepath}')
print(f'[+] {hashtype.upper()}: {hash_obj.hexdigest()}')
