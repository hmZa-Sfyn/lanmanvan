#!/usr/bin/env python3
"""
Hash Generator Module
Generates various hash values for input
"""

import hashlib
import os

def generate_hashes(data):
    """Generate multiple hash values"""
    hashes = {
        'md5': hashlib.md5(data.encode()).hexdigest(),
        'sha1': hashlib.sha1(data.encode()).hexdigest(),
        'sha256': hashlib.sha256(data.encode()).hexdigest(),
        'sha512': hashlib.sha512(data.encode()).hexdigest(),
    }
    return hashes

def main():
    data = os.getenv('ARG_DATA') or os.getenv('ARG_INPUT') or 'test'
    
    print(f"[*] Generating hashes for: {data}\n")
    
    hashes = generate_hashes(data)
    
    for hash_type, hash_value in hashes.items():
        print(f"  {hash_type.upper():8} : {hash_value}")
    
    print()

if __name__ == '__main__':
    main()
