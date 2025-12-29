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
        lat = data.get("lat")
        lon = data.get("lon")
        maps_url = f'https://www.google.com/maps/search/{lat},{lon}'
        
        print(f'[+] IP: {data.get("query")}')
        print(f'[+] Country: {data.get("country")}')
        print(f'[+] Region: {data.get("regionName")}')
        print(f'[+] City: {data.get("city")}')
        print(f'[+] ZIP: {data.get("zip")}')
        print(f'[+] Coordinates: {lat}, {lon}')
        print(f'[+] Google Maps: {maps_url}')
        print(f'[+] ISP: {data.get("isp")}')
        print(f'[+] Organization: {data.get("org")}')
        print(f'[+] AS: {data.get("as")}')
        print(f'[+] Timezone: {data.get("timezone")}')
        print(f'[+] Mobile: {data.get("mobile")}')
        print(f'[+] Proxy: {data.get("proxy")}')
        print(f'[+] Hosting: {data.get("hosting")}')
    else:
        print('[!] Could not resolve IP')
except Exception as e:
    print(f'[!] Error: {e}')
    try:
        r2 = requests.get(f'https://ipwhois.app/json/{ip}', timeout=5)
        data2 = r2.json()
        if data2.get('success'):
            print(f'[+] Continent: {data2.get("continent")}')
            print(f'[+] Type: {data2.get("type")}')
            print(f'[+] Network: {data2.get("isp")}')
    except:
        pass
    
    try:
        r3 = requests.get(f'https://api.abuseipdb.com/api/v2/check', 
                         headers={'Key': os.getenv('ABUSEIPDB_KEY', ''), 'Accept': 'application/json'},
                         params={'ipAddress': ip, 'maxAgeInDays': '90'}, timeout=5)
        data3 = r3.json()
        if data3.get('data'):
            print(f'[+] Abuse Score: {data3["data"].get("abuseConfidenceScore")}%')
            print(f'[+] Total Reports: {data3["data"].get("totalReports")}')
    except:
        pass
    
    try:
        r4 = requests.get(f'https://api.ip2location.io/', 
                         params={'ip': ip, 'key': os.getenv('IP2LOCATION_KEY', '')}, timeout=5)
        data4 = r4.json()
        if data4.get('country_code'):
            print(f'[+] Usage Type: {data4.get("usage_type")}')
            print(f'[+] Threat Level: {data4.get("threat")}')
    except:
        pass