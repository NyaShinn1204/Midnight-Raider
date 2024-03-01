import os

import bypass.header as header
from colorama import Fore

pretty = Fore.LIGHTMAGENTA_EX + Fore.LIGHTCYAN_EX

def get_filename():
  return os.path.basename(__file__)    

def printl(num, data):
  if num == "error":
    print(f"[{Fore.LIGHTRED_EX}Error{Fore.RESET}] [{get_filename()}] " + data)
  if num == "info":
    print(f"[{Fore.LIGHTGREEN_EX}Info{Fore.RESET}] [{get_filename()}] " + data)

def get_fingerprint():
  session = header.get_session.get_session()
  req_header = header.request_header(None)
  headers = req_header
  response = session.get('https://discord.com/api/v9/experiments', headers=headers)
  if response.status_code == 200:
    data = response.json()
    fingerprint = data["fingerprint"]
    printl("info", f"{pretty}Got Fingerprint {Fore.RESET}| " + Fore.GREEN + fingerprint + Fore.RESET)
    return fingerprint
  else:
    printl("error", f"{pretty}Failed Got Fingerprint {Fore.RESET}| " + Fore.RESET)
    return