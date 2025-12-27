#!/usr/bin/env python3
import os, sys, requests
try:
    from bs4 import BeautifulSoup
except:
    import subprocess
    subprocess.run([sys.executable, '-m', 'pip', 'install', 'beautifulsoup4', '-q'])
    from bs4 import BeautifulSoup

url = os.getenv('ARG_URL', '')
selector = os.getenv('ARG_SELECTOR', 'p')
if not url:
    print('[!] URL required')
    sys.exit(1)
try:
    r = requests.get(url, timeout=10)
    soup = BeautifulSoup(r.content, 'html.parser')
    elements = soup.select(selector)
    print(f'[+] Found {len(elements)} elements')
    for i, elem in enumerate(elements[:10]):
        print(f'{i+1}. {elem.get_text()[:100]}')
except Exception as e:
    print(f'[!] Error: {e}')