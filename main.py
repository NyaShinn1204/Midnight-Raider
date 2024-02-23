import os
import time
import json
import threading
import webbrowser
import subprocess
import requests
import colorama
import tkinter as tk
import customtkinter as ctk
from PIL import Image
from pystyle import *
from colorama import Fore
from customtkinter import *
from CTkMessagebox import CTkMessagebox
from CTkToolTip import *

import rpc.example as rpc

version = "1.0.0"

def printl(num, data):
  filename = os.path.basename(__file__)
  if num == "error":
    print(f"["+Colorate.Horizontal(Colors.red_to_blue, "Error")+"]"+f"[{filename}] " + data)
    #print(f"[{Fore.LIGHTRED_EX}Error{Fore.RESET}] [{filename}] " + data)
  if num == "debug":
    print(f"["+Colorate.Horizontal(Colors.cyan_to_blue, "Debug")+"]"+f"[{filename}] " + data)
    #print(f"[{Fore.LIGHTCYAN_EX}Debug{Fore.RESET}] [{filename}] " + data)
  if num == "info":
    print(f"["+Colorate.Horizontal(Colors.white_to_blue, "Info")+"]"+f"[{filename}] " + data)
    #print(f"[{Fore.LIGHTGREEN_EX}Info{Fore.RESET}] [{filename}] " + data)
    
def get_hwid():
  try:
    if os.name == 'posix':
      cmd = 'cat /etc/machine-id'
      uuid = subprocess.check_output(cmd,shell=True)
      uuid = uuid[:-1].decode('utf-8')
      return uuid
    if os.name == 'nt':
      cmd = 'powershell -Command (Get-WmiObject -Class Win32_ComputerSystemProduct).UUID'
      uuid = subprocess.run(cmd, capture_output=True, text=True, shell=True)
      uuid = uuid.stdout.strip()
      return uuid
  except:
    printl("error", "get_hwid exception error wrong")

System.Size(120, 30)
System.Clear()

root = tk.Tk()
root.geometry("1280x720")
root.resizable(0, 0)
root.title("Midnight Raider | "+version)
root.iconbitmap("./data/icon.ico")
root.configure(bg="#baebfb")

def load_background():
  ctk.CTkLabel(master=root,image=ctk.CTkImage(Image.open("./data/background-01.jpg"),size=(1280,720)),text="").pack()

# Load First
logo = f"""
                    ___  ____     _       _       _     _    ______      _     _           
    ..&@.           |  \/  (_)   | |     (_)     | |   | |   | ___ \    (_)   | |          
  ..@@@&.           | .  . |_  __| |_ __  _  __ _| |__ | |_  | |_/ /__ _ _  __| | ___ _ __ 
  .&&&&&,..         | |\/| | |/ _` | '_ \| |/ _` | '_ \| __| |    // _` | |/ _` |/ _ \ '__|
  ..&&&&&&&#.       | |  | | | (_| | | | | | (_| | | | | |_  | |\ \ (_| | | (_| |  __/ |   
    ..#&&&...       \_|  |_/_|\__,_|_| |_|_|\__, |_| |_|\__| \_| \_\__,_|_|\__,_|\___|_|   
                                             __/ |                                         
                                            |___/                                      
"""
info = f"""
HWID: [{get_hwid()}]                Version: [{version}]
"""
print(Colorate.Horizontal(Colors.white_to_blue, Center.XCenter(logo)))
print(Colorate.Color(Colors.white, Center.XCenter(info)))
print("""\n------------------------------------------------------------------------------------------------------------------------""")
printl("info", "Loading GUI")
load_background()
printl("info", "Starting RPC")
threading.Thread(target=rpc.start).start()


root.mainloop()