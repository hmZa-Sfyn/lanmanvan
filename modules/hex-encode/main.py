#!/usr/bin/env python3
import os, sys
data = os.getenv('ARG_DATA', '')
print('[+] Original:', data)
print('[+] Hex:', data.encode().hex())
