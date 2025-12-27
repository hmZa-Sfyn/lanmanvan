#!/usr/bin/env python3
import os, sys
try:
    from scapy.all import sniff, IP, ICMP
except:
    print('[*] Installing scapy...')
    import subprocess
    subprocess.run([sys.executable, '-m', 'pip', 'install', 'scapy', '-q'])
    from scapy.all import sniff, IP, ICMP

interface = os.getenv('ARG_INTERFACE', '')
count = int(os.getenv('ARG_COUNT', '10'))

def packet_callback(packet):
    if IP in packet:
        ip_src = packet[IP].src
        ip_dst = packet[IP].dst
        print(f'[+] {ip_src} -> {ip_dst}')
        if ICMP in packet:
            print(f'    ICMP Request')

print(f'[*] Sniffing {count} packets...')
try:
    if interface:
        sniff(prn=packet_callback, iface=interface, count=count)
    else:
        sniff(prn=packet_callback, count=count)
except PermissionError:
    print('[!] Packet sniffing requires root/admin privileges')
except Exception as e:
    print(f'[!] Error: {e}')
