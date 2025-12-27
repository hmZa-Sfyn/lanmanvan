#!/usr/bin/env python3
import os, sys, socket, threading
domain = os.getenv('ARG_DOMAIN', '')
wordlist = os.getenv('ARG_WORDLIST', 'www,mail,ftp,localhost,webmail,smtp,pop,ns,admin,test,portal,api,dev,staging').split(',')

if not domain:
    print('[!] Domain required')
    sys.exit(1)

found = []
def check_subdomain(sub):
    try:
        full = f'{sub}.{domain}'
        ip = socket.gethostbyname(full)
        found.append((full, ip))
        print(f'[+] Found: {full} -> {ip}')
    except:
        pass

print(f'[*] Enumerating subdomains for {domain}...')
threads = [threading.Thread(target=check_subdomain, args=(sub.strip(),)) for sub in wordlist]
for t in threads:
    t.start()
for t in threads:
    t.join()

print(f'[+] Found {len(found)} subdomains')
