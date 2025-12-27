#!/usr/bin/env python3
import os, sys, socket
domain = os.getenv('ARG_DOMAIN', '')
if not domain:
    print('[!] Domain required')
    sys.exit(1)
try:
    ip = socket.gethostbyname(domain)
    print(f'[+] Domain: {domain}')
    print(f'[+] IP: {ip}')
    # Try reverse DNS
    try:
        host = socket.gethostbyaddr(ip)
        print(f'[+] Reverse: {host[0]}')
    except:
        pass
except socket.gaierror:
    print('[!] Could not resolve domain')
