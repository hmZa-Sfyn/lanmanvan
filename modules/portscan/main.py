#!/usr/bin/env python3
"""
Port Scanner Module
Scans ports on a target host
"""

import sys
import socket
import os

def scan_port(host, port, timeout=1):
    """Attempt to connect to a port"""
    try:
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        result = sock.connect_ex((host, port))
        sock.close()
        return result == 0
    except socket.error:
        return False

def main():
    # Get arguments from environment variables
    host = os.getenv('ARG_HOST') or os.getenv('ARG_TARGET') or 'localhost'
    port_range = os.getenv('ARG_PORTS') or '80,443,22,21,3306,5432'
    
    ports = []
    for p in port_range.split(','):
        try:
            ports.append(int(p.strip()))
        except:
            pass
    
    if not ports:
        ports = [80, 443, 22, 21, 3306, 5432]
    
    print(f"[*] Scanning {host} for open ports...")
    print(f"[*] Ports to scan: {ports}\n")
    
    open_ports = []
    for port in ports:
        if scan_port(host, port):
            open_ports.append(port)
            print(f"[+] Port {port} is OPEN")
        else:
            print(f"[-] Port {port} is closed")
    
    print(f"\n[*] Scan complete. Found {len(open_ports)} open ports")
    if open_ports:
        print(f"[+] Open ports: {open_ports}")

if __name__ == '__main__':
    main()
