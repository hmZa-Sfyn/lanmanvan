#!/usr/bin/env python3
"""
Host Discovery Module - Discover live hosts on network
Demonstrates: ICMP ping, threading, network scanning
"""

import os
import sys
import subprocess
import concurrent.futures
from ipaddress import ip_network, ip_address
import time

def ping_host(host, timeout=1):
    """Ping a host and return True if alive"""
    try:
        result = subprocess.run(
            ["ping", "-c", "1", "-W", str(timeout), str(host)],
            capture_output=True,
            timeout=timeout + 1
        )
        return result.returncode == 0
    except:
        return False

def main():
    network = os.getenv('ARG_NETWORK') or '192.168.1.0/24'
    threads = int(os.getenv('ARG_THREADS') or '10')
    
    print(f"[*] Starting host discovery on {network}")
    print(f"[*] Using {threads} threads")
    print()
    
    try:
        net = ip_network(network, strict=False)
    except ValueError:
        print("[!] Invalid network address")
        sys.exit(1)
    
    alive_hosts = []
    total_hosts = len(list(net.hosts()))
    scanned = 0
    
    print(f"[*] Total hosts to scan: {total_hosts}")
    print()
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=threads) as executor:
        futures = {executor.submit(ping_host, host): host for host in net.hosts()}
        
        for future in concurrent.futures.as_completed(futures):
            host = futures[future]
            scanned += 1
            
            try:
                is_alive = future.result()
                if is_alive:
                    alive_hosts.append(str(host))
                    print(f"[+] Host {host} is ALIVE")
                else:
                    print(f"[-] Host {host} is DOWN")
            except Exception as e:
                print(f"[!] Error scanning {host}: {e}")
            
            # Progress indicator
            if scanned % 10 == 0:
                progress = (scanned / total_hosts) * 100
                print(f"[*] Progress: {scanned}/{total_hosts} ({progress:.1f}%)")
    
    print()
    print(f"[*] Scan complete!")
    print(f"[+] Found {len(alive_hosts)} alive hosts:")
    for host in sorted(alive_hosts):
        print(f"    {host}")

if __name__ == '__main__':
    main()
