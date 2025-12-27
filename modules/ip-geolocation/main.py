#!/usr/bin/env python3
import os, sys, requests
ip = os.getenv('ARG_IP', '')
if not ip:
    print('[!] IP required')
    sys.exit(1)
try:
    r = requests.get(f'http://ip-api.com/json/{ip}', timeout=5)
    data = r.json()
    if data['status'] == 'success':
        print(f'[+] IP: {data.get("query")}')
        print(f'[+] Country: {data.get("country")}')
        print(f'[+] City: {data.get("city")}')
        print(f'[+] Lat/Lon: {data.get("lat")}, {data.get("lon")}')
        print(f'[+] ISP: {data.get("isp")}')
    else:
        print('[!] Could not resolve IP')
except Exception as e:
    print(f'[!] Error: {e}')
