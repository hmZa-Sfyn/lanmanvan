#!/usr/bin/env python3
import os
import sys
import requests
import json
import socket
import subprocess
from datetime import datetime

# Optional: pip install scapy (for probing, requires root privileges)
try:
    from scapy.all import IP, ICMP, TCP, UDP, sr1, RandShort
    SCAPY_AVAILABLE = True
except ImportError:
    SCAPY_AVAILABLE = False

def print_info(label, value):
    if value is not None and value != "" and value != False:
        print(f'[+] {label}: {value}')

ip = os.getenv('ARG_IP', '')
scan_ports = os.getenv('ARG_SCAN_PORTS', 'false').lower() == 'true'
api_keys = {}  # Populate with os.getenv if needed, e.g., api_keys['abuseipdb'] = os.getenv('ABUSEIPDB_KEY')

if not ip:
    print('[!] IP required')
    sys.exit(1)

headers = {'User-Agent': 'ultimate-ip-tool/1.2'}

try:
    # Reverse DNS
    print('[*] Performing reverse DNS lookup...')
    try:
        hostname, _, _ = socket.gethostbyaddr(ip)
        print_info('Hostname', hostname)
    except socket.herror:
        print('[!] No reverse DNS record')

    # Primary: ip-api.com
    print('\n[*] Querying ip-api.com...')
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

    # AbuseIPDB
    print('\n[*] Querying AbuseIPDB...')
    abuse_key = api_keys.get('abuseipdb', '')
    abuse_headers = {'Key': abuse_key, 'Accept': 'application/json'} if abuse_key else {'Accept': 'application/json'}
    abuse_url = f'https://api.abuseipdb.com/api/v2/check?ipAddress={ip}&maxAgeInDays=90&verbose'
    r_abuse = requests.get(abuse_url, headers=abuse_headers, timeout=10)
    if r_abuse.status_code == 200:
        abuse_data = r_abuse.json()['data']
        print_info('Abuse Confidence Score', f'{abuse_data.get("abuseConfidenceScore")}%')
        print_info('Total Reports', abuse_data.get('totalReports'))
        print_info('Last Reported', abuse_data.get('lastReportedAt'))
        print_info('Is Whitelisted', abuse_data.get('isWhitelisted'))
        print_info('Is Tor', abuse_data.get('isTor'))
        print_info('Usage Type', abuse_data.get('usageType'))
        if 'reports' in abuse_data and abuse_data['reports']:
            print('[+] Recent Abuse Reports:')
            for report in abuse_data['reports'][:5]:  # Limit to 5
                print(f'  - {report.get("reportedAt")}: {report.get("comment")}')
    else:
        print(f'[!] AbuseIPDB error: {r_abuse.status_code} - Sign up for free key at abuseipdb.com for full access')

    # ipwhois.app
    print('\n[*] Querying ipwhois.app...')
    r_whois = requests.get(f'http://ipwhois.app/json/{ip}', timeout=10, headers=headers)
    if r_whois.status_code == 200:
        whois_data = r_whois.json()
        if whois_data.get('success'):
            print_info('Connection Type', whois_data.get('type'))
            print_info('WHOIS ISP', whois_data.get('isp'))
            print_info('WHOIS Org', whois_data.get('org'))
            print_info('Abuse Email', whois_data.get('abuse_email'))

    # Shodan for OS, vulns, etc. (requires API key)
    print('\n[*] Querying Shodan (requires API key)...')
    shodan_key = api_keys.get('shodan', '')
    if shodan_key:
        r_shodan = requests.get(f'https://api.shodan.io/shodan/host/{ip}?key={shodan_key}', timeout=10)
        if r_shodan.status_code == 200:
            shodan_data = r_shodan.json()
            print_info('OS (Shodan)', shodan_data.get('os'))
            print_info('Hostnames', ', '.join(shodan_data.get('hostnames', [])))
            print_info('Ports', ', '.join(map(str, shodan_data.get('ports', []))))
            if 'data' in shodan_data:
                print('[+] Open Services:')
                for service in shodan_data['data'][:5]:
                    print(f'  - Port {service.get("port")}: {service.get("product", "Unknown")} {service.get("version", "")}')
                    if 'vulns' in service:
                        print(f'    Vulns: {", ".join(service["vulns"])}')
    else:
        print('[!] Shodan skipped: Get free API key at shodan.io')

    # Traceroute (using subprocess)
    print('\n[*] Performing traceroute...')
    try:
        trace_out = subprocess.check_output(['traceroute', '-m', '20', ip], timeout=30).decode()
        print('[+] Traceroute:')
        print(trace_out)
    except Exception as e:
        print(f'[!] Traceroute failed: {e} (ensure traceroute installed)')

    # Probing (ICMP, TCP, UDP) - requires scapy and root
    if scan_ports:
        if not SCAPY_AVAILABLE:
            print('[!] Scapy not available: pip install scapy')
        else:
            print('\n[*] Probing (requires root)...')
            # ICMP Ping
            ping_pkt = IP(dst=ip)/ICMP()
            ping_resp = sr1(ping_pkt, timeout=2, verbose=0)
            if ping_resp:
                print_info('ICMP Ping', 'Alive')
                print_info('TTL (for OS guess)', ping_resp.ttl)  # e.g., 64=Linux, 128=Windows, 255=Solaris
                if ping_resp.ttl <= 64:
                    os_guess = 'Likely Linux/Unix'
                elif ping_resp.ttl <= 128:
                    os_guess = 'Likely Windows'
                else:
                    os_guess = 'Likely Solaris/Network Device'
                print_info('OS Guess (from TTL)', os_guess)
            else:
                print('[!] ICMP Ping: No response')

            # Common ports to scan
            common_ports = [22, 80, 443, 3389, 21, 25, 53, 110, 143, 445]

            print('[+] TCP SYN Scan:')
            for port in common_ports:
                tcp_pkt = IP(dst=ip)/TCP(sport=RandShort(), dport=port, flags='S')
                tcp_resp = sr1(tcp_pkt, timeout=2, verbose=0)
                if tcp_resp and tcp_resp.haslayer(TCP):
                    if tcp_resp[TCP].flags == 0x12:  # SYN-ACK
                        print(f'  - Port {port}: Open')
                    elif tcp_resp[TCP].flags == 0x14:  # RST
                        print(f'  - Port {port}: Closed')
                else:
                    print(f'  - Port {port}: Filtered/No response')

            print('[+] UDP Probe (common ports):')
            for port in [53, 123, 161]:  # DNS, NTP, SNMP
                udp_pkt = IP(dst=ip)/UDP(sport=RandShort(), dport=port)
                udp_resp = sr1(udp_pkt, timeout=2, verbose=0)
                if udp_resp:
                    print(f'  - Port {port}: Response received (possibly open)')
                else:
                    print(f'  - Port {port}: No response (possibly open/filtered)')

except requests.Timeout:
    print('[!] Request timeout')
except requests.RequestException as e:
    print(f'[!] Network error: {e}')
except Exception as e:
    print(f'[!] Error: {e}')