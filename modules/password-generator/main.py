#!/usr/bin/env python3
import os, sys, random, string
length = int(os.getenv('ARG_LENGTH', '16'))
count = int(os.getenv('ARG_COUNT', '1'))
chars = string.ascii_letters + string.digits + '!@#$%^&*()'
for i in range(count):
    pwd = ''.join(random.choice(chars) for _ in range(length))
    print(f'[+] Password {i+1}: {pwd}')
