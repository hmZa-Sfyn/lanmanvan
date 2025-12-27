#!/usr/bin/env python3
import os, sys, ssl, socket
from datetime import datetime
domain = os.getenv('ARG_DOMAIN', '')
if not domain:
    print('[!] Domain required')
    sys.exit(1)
try:
    context = ssl.create_default_context()
    with socket.create_connection((domain, 443), timeout=5) as sock:
        with context.wrap_socket(sock, server_hostname=domain) as ssock:
            cert = ssock.getpeercert()
            print(f'[+] Domain: {domain}')
            print(f'[+] Subject: {cert.get("subject", "N/A")}')
            print(f'[+] Issuer: {cert.get("issuer", "N/A")}')
            print(f'[+] Version: {cert.get("version", "N/A")}')
except Exception as e:
    print(f'[!] Error: {e}')