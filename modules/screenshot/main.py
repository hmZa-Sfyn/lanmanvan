#!/usr/bin/env python3
import os, sys
try:
    from selenium import webdriver
except:
    print('[*] Installing selenium...')
    import subprocess
    subprocess.run([sys.executable, '-m', 'pip', 'install', 'selenium', '-q'])
    from selenium import webdriver

url = os.getenv('ARG_URL', '')
output = os.getenv('ARG_OUTPUT', '/tmp/screenshot.png')

if not url:
    print('[!] URL required')
    sys.exit(1)

try:
    print(f'[*] Capturing screenshot of {url}')
    options = webdriver.ChromeOptions()
    options.add_argument('--headless')
    options.add_argument('--no-sandbox')
    options.add_argument('--disable-dev-shm-usage')
    driver = webdriver.Chrome(options=options)
    driver.get(url)
    driver.save_screenshot(output)
    driver.quit()
    print(f'[+] Screenshot saved to {output}')
except Exception as e:
    print(f'[!] Error: {e}')
    print('[!] Make sure Chrome and ChromeDriver are installed')
