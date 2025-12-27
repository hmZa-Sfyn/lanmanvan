#!/usr/bin/env python3
import os, sys, socket
port = int(os.getenv('ARG_PORT', '0'))
if port < 1 or port > 65535:
    print('[!] Invalid port')
    sys.exit(1)

common_ports = {
    20: 'FTP Data', 21: 'FTP', 22: 'SSH', 23: 'Telnet', 25: 'SMTP',
    53: 'DNS', 80: 'HTTP', 110: 'POP3', 143: 'IMAP', 443: 'HTTPS',
    445: 'SMB', 3306: 'MySQL', 3389: 'RDP', 5432: 'PostgreSQL',
    5900: 'VNC', 6379: 'Redis', 8080: 'HTTP Alt', 8443: 'HTTPS Alt',
    9200: 'Elasticsearch', 27017: 'MongoDB'
}

try:
    service = socket.getservbyport(port)
except:
    service = common_ports.get(port, 'Unknown')

print(f'[+] Port: {port}')
print(f'[+] Service: {service}')
print(f'[+] Protocol: TCP/UDP')
