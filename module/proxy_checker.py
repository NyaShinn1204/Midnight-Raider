import random
import socket
import traceback
import urllib.request
from threading import Thread

socket.setdefaulttimeout(30)
    
def check(update_proxy, proxies, types):
    threads = []
    for proxy in proxies:
        thread = Thread(target=check_proxy, args=(update_proxy, proxy,types))
        thread.start()
        threads.append(thread)
    for thread in threads:
        thread.join()
        
def check_proxy(update_proxy, pip, types):
    try:
        proxy_handler = urllib.request.ProxyHandler({types: pip})        
        opener = urllib.request.build_opener(proxy_handler)
        opener.addheaders = [('User-agent', 'Mozilla/5.0')]
        urllib.request.install_opener(opener)        
        sock=urllib.request.urlopen('http://www.google.com')
    except urllib.error.HTTPError as e:
        update_proxy(False, pip)
        print("[-] Not working: " +pip)
        return
    except Exception as detail:
        update_proxy(False, pip)
        print("[-] Not working: " +pip)
        return
    print("[+] Working: " + pip)
    update_proxy(True, pip)
    