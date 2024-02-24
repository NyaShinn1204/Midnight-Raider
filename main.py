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


def module_list_frame():
  #global modulelist
  tk.Label(root, bg="#baebfb", width=1024, height=720).place(x=0,y=0)
  tk.Label(root, bg="#0f1314", width=32, height=720).place(x=0,y=0)
  ctk.CTkLabel(master=root,image=ctk.CTkImage(Image.open("./data/icon.png"),size=(85,85)),text="").place(x=5,y=5)
  tk.Label(root, bg="#0f1314", text="Midnight", fg="#fff", font=("Dubai Medium", 20)).place(x=100,y=0)
  tk.Label(root, bg="#0f1314", text="Raider", fg="#fff", font=("Dubai Medium", 20)).place(x=100,y=40)
  
#  modulelist = ctk.CTkFrame(master=root, width=230, height=720, corner_radius=0, fg_color=c4)
#  modulelist.place(x=0,y=100)
#  tk.Canvas(bg=c6, highlightthickness=0, height=2080, width=4).place(x=230, y=0)
#  ctk.CTkButton(master=modulelist, image=ctk.CTkImage(Image.open("data/join_leave.png"),size=(20, 20)), compound="left", fg_color=c4, hover_color=c5, corner_radius=0, text=lang_load_set("joiner_leaver"), width=195, height=40, font=set_fonts(16, "bold"), anchor="w", command= lambda: module_scroll_frame(1, 1)).place(x=20,y=12)
#  ctk.CTkButton(master=modulelist, image=ctk.CTkImage(Image.open("data/spammer.png"),size=(20, 20)), compound="left", fg_color=c4, hover_color=c5, corner_radius=0, text=lang_load_set("spammer"), width=195, height=40, font=set_fonts(16, "bold"), anchor="w", command= lambda: module_scroll_frame(1, 2)).place(x=20,y=57)
#  ctk.CTkButton(master=modulelist, image=ctk.CTkImage(Image.open("data/setting.png"),size=(20, 20)), compound="left", fg_color=c4, hover_color=c5, corner_radius=0, text=lang_load_set("settings"), width=195, height=40, font=set_fonts(16, "bold"), anchor="w", command= lambda: module_scroll_frame(2, 1)).place(x=20,y=516)
#  ctk.CTkButton(master=modulelist, image=ctk.CTkImage(Image.open("data/info.png"),size=(20, 20)), compound="left", fg_color=c4, hover_color=c5, corner_radius=0, text=lang_load_set("about"), width=195, height=40, font=set_fonts(16, "bold"), anchor="w", command= lambda: module_scroll_frame(2, 2)).place(x=20,y=562)
#  
#  credit_frame = ctk.CTkFrame(root, width=1020, height=50, fg_color=c1, bg_color=c2)
#  credit_frame.place(x=245, y=10)
#  ctk.CTkButton(master=credit_frame, image=ctk.CTkImage(Image.open("data/link.png"),size=(20, 20)), compound="right", fg_color=c1, text_color="#fff", corner_radius=0, text="", width=20, height=20, font=set_fonts(16, None), anchor="w", command= lambda: CTkMessagebox(title="Version Info", message=f"Version: {version}\n\nDeveloper: {developer}\nTester: {testers}", width=450)).place(x=10,y=10)
#  ctk.CTkLabel(master=credit_frame, fg_color=c1, text_color="#fff", corner_radius=0, text=lang_load_set("username")+": "+os.getlogin(), width=20, height=20, font=set_fonts(16, "bold"), anchor="w").place(x=40,y=5)
#  ctk.CTkLabel(master=credit_frame, fg_color=c1, text_color="#fff", corner_radius=0, text="Hwid: "+get_hwid(), width=20, height=20, font=set_fonts(16, "bold"), anchor="w").place(x=40,y=25)


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
#load_background()
module_list_frame()
printl("info", "Starting RPC")
threading.Thread(target=rpc.start).start()

0
root.mainloop()