#!/usr/bin/env python3
import os
import sys
import requests
import json

def print_info(label, value):
    if value is not None and value != "" and value != False:
        print(f'[+] {label}: {value}')

ip = os.getenv('ARG_IP', '')
if not ip:
    print('[!] IP required')
    sys.exit(1)

headers = {'User-Agent': 'advanced-ip-tool/1.1'}

try:
    # Primary: ip-api.com (no key, good fields, 45 req/min limit)
    print('[*] Querying ip-api.com...')
    r1 = requests.get(f'http://ip-api.com/json/{ip}?fields=status,message,continent,continentCode,country,countryCode,region,regionName,city,district,zip,lat,lon,timezone,offset,currency,isp,org,as,asname,mobile,proxy,hosting,query', timeout=10, headers=headers)
    data1 = r1.json()
    if data1.get('status') != 'success':
        raise Exception(data1.get('message', 'Failed'))

    print(f'[+] IP: {data1.get("query")}')
    print_info('Continent', f'{data1.get("continent")} ({data1.get("continentCode")})')
    print_info('Country', f'{data1.get("country")} ({data1.get("countryCode")})')
    print_info('Region', f'{data1.get("regionName")} ({data1.get("region")})')
    print_info('City', data1.get('city'))
    print_info('District', data1.get('district'))
    print_info('ZIP', data1.get('zip'))
    print_info('Coordinates', f'{data1.get("lat")}, {data1.get("lon")}')
    print(f'[+] Google Maps: https://www.google.com/maps/search/{data1.get("lat")},{data1.get("lon")}')
    print_info('Timezone', f'{data1.get("timezone")} (offset: {data1.get("offset")})')
    print_info('Currency', data1.get('currency'))
    print_info('ISP', data1.get('isp'))
    print_info('Organization', data1.get('org'))
    print_info('AS', f'{data1.get("as")} ({data1.get("asname")})')
    print_info('Mobile', data1.get('mobile'))
    print_info('Proxy/VPN', data1.get('proxy'))
    print_info('Hosting', data1.get('hosting'))

    # Additional: AbuseIPDB for threat/abuse reports (free, requires key but works without for basic)
    print('\n[*] Querying AbuseIPDB for abuse reports...')
    abuse_url = f'https://api.abuseipdb.com/api/v2/check?ipAddress={ip}&maxAgeInDays=90&verbose'
    r_abuse = requests.get(abuse_url, headers={'Key': '', 'Accept': 'application/json'}, timeout=10)  # Add your free key if you have one
    if r_abuse.status_code == 200:
        abuse_data = r_abuse.json()['data']
        print_info('Abuse Confidence Score', f'{abuse_data.get("abuseConfidenceScore")}%')
        print_info('Total Reports', abuse_data.get('totalReports'))
        print_info('Last Reported', abuse_data.get('lastReportedAt'))
        print_info('Is Whitelisted', abuse_data.get('isWhitelisted'))
        print_info('Is Tor', abuse_data.get('isTor'))
        usage_type = abuse_data.get('usageType')
        if usage_type:
            print_info('Usage Type', usage_type)
    else:
        print('[!] AbuseIPDB: Rate limited or no key (sign up for free key at abuseipdb.com for reports)')

    # Optional extra: whois via ipwhois.io (includes some connection type)
    print('\n[*] Querying ipwhois.io for additional WHOIS/connection info...')
    r_whois = requests.get(f'http://ipwhois.app/json/{ip}', timeout=10, headers=headers)
    if r_whois.status_code == 200:
        whois_data = r_whois.json()
        if whois_data.get('success'):
            print_info('Connection Type', whois_data.get('type'))  # e.g., Residential, Business
            print_info('WHOIS ISP', whois_data.get('isp'))
            print_info('WHOIS Org', whois_data.get('org'))
    else:
        print('[!] ipwhois failed')

except requests.Timeout:
    print('[!] Request timeout')
except requests.RequestException as e:
    print(f'[!] Network error: {e}')
except Exception as e:
    print(f'[!] Error: {e}')