#!/usr/bin/env python3
import os, sys, re
email = os.getenv('ARG_EMAIL', '')
if not email:
    print('[!] Email required')
    sys.exit(1)

pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
if re.match(pattern, email):
    print(f'[+] Valid email: {email}')
else:
    print(f'[!] Invalid email: {email}')
