#!/usr/bin/env python3
import os, sys
data = os.getenv('ARG_DATA', '')
try:
    decoded = bytes.fromhex(data).decode()
    print('[+] Hex:', data)
    print('[+] Decoded:', decoded)
except:
    print('[!] Invalid hex data')
