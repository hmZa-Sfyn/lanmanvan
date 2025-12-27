#!/usr/bin/env python3
import os, sys, requests, threading
url = os.getenv('ARG_URL', '').rstrip('/')
wordlist = os.getenv('ARG_WORDLIST', 'admin,backup,config,wp-admin,api,test,debug,upload,files,doc,public').split(',')
if not url:
    print('[!] URL required')
    sys.exit(1)
found = []
def check_dir(d):
    try:
        r = requests.head(f'{url}/{d.strip()}', timeout=3)
        if r.status_code < 400:
            found.append(d.strip())
            print(f'[+] {d.strip()} ({r.status_code})')
    except:
        pass
print(f'[*] Scanning {url}...')
threads = [threading.Thread(target=check_dir, args=(d,)) for d in wordlist]
for t in threads: t.start()
for t in threads: t.join()
print(f'[+] Found {len(found)} directories')