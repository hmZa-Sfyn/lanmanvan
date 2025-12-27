#!/usr/bin/env python3
import os, sys
try:
    import whois
except:
    import subprocess
    subprocess.run([sys.executable, '-m', 'pip', 'install', 'python-whois', '-q'])
    import whois

domain = os.getenv('ARG_DOMAIN', '')
if not domain:
    print('[!] Domain required')
    sys.exit(1)

try:
    w = whois.whois(domain)
    print(f'[+] Domain: {w.domain}')
    print(f'[+] Registrar: {w.registrar}')
    print(f'[+] Created: {w.creation_date}')
    print(f'[+] Expires: {w.expiration_date}')
    if hasattr(w, 'name_servers'):
        print(f'[+] Name Servers: {w.name_servers}')
except Exception as e:
    print(f'[!] Error: {e}')
