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
developer = "NyaShinn1204"
contributors = "None"
testers = "None"

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

def get_plan():
  try:
    if get_hwid() == "BFA3C0C0-77D2-5F57-D558-663A669C1043":
      return "Developer"
    else:
      return "Normal"
  except:
    printl("error", "get_plan exception error wrong")

def gui_close():
  root.destroy()
  rpc.stop_threads = True

System.Size(120, 30)
System.Clear()

root = tk.Tk()
root.geometry("1280x720")
root.resizable(0, 0)
root.title("Midnight Raider | "+version)
root.iconbitmap("./data/icon.ico")
root.configure(bg="#baebfb")
root.protocol("WM_DELETE_WINDOW", gui_close)

def load_background():
  ctk.CTkLabel(master=root,image=ctk.CTkImage(Image.open("./data/background-01.jpg"),size=(1280,720)),text="").pack()

def clear_frame(frame):
  for widget in frame.winfo_children():
    widget.destroy()
  frame.pack_forget()

def module_scroll_frame(num1, num2):
  print("a")
  global module_frame
  frame_scroll = module_frame = ctk.CTkScrollableFrame(root, fg_color="#0f1314", bg_color="#152945", width=1000, height=630)
  module_frame.place(x=245, y=70)
  clear_frame(frame_scroll)
  if num2 == 2:
    if num2 == 2:
      # About
      credits_frame = ctk.CTkFrame(module_frame, width=940, height=400, border_width=0, corner_radius=5, fg_color="#0f1314")
      credits_frame.grid(row=1, column=1, padx=6, pady=6)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text="Midnight Raider github:", font=("Roboto", 12)).place(x=0,y=0)
      test = tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text="Github link", font=("Roboto", 12, "underline"))
      test.place(x=175,y=0)
      test.bind("<Button-1>", lambda e:webbrowser.open_new("https://github.com/NyaShinn1204/Midnight-Raider"))
      tk.Label(credits_frame, bg="#0f1314", fg="#fff", text="Main developer and updater:", font=("Roboto", 12)).place(x=0,y=25)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text=developer, font=("Roboto", 12)).place(x=210,y=25)
      tk.Label(credits_frame, bg="#0f1314", fg="#fff", text="Main contributors:", font=("Roboto", 12)).place(x=0,y=50)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text=contributors, font=("Roboto", 12)).place(x=137,y=50)
      tk.Label(credits_frame, bg="#0f1314", fg="#fff", text="Main testers:", font=("Roboto", 12)).place(x=0,y=75)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text=testers, font=("Roboto", 12)).place(x=100,y=75)
      
      tk.Label(credits_frame, bg="#0f1314", fg="#fff", text="Respect:", font=("Roboto", 12)).place(x=0,y=125)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text="Akebi GC", font=("Roboto", 12)).place(x=15,y=145)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text="Bkebi GC", font=("Roboto", 12)).place(x=15,y=165)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text="TwoCoinRaider", font=("Roboto", 12)).place(x=15,y=185)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text="ThreeCoinRaider", font=("Roboto", 12)).place(x=15,y=206)
      tk.Label(credits_frame, bg="#0f1314", fg="#c9f7fe", text="RaizouRaider", font=("Roboto", 12)).place(x=15,y=227)

      printl("debug", "Open About Tab")

def module_list_frame():
  #global modulelist
  tk.Label(root, bg="#152945", width=1024, height=720).place(x=0,y=0)
  tk.Label(root, bg="#0f1314", width=32, height=720).place(x=0,y=0)
  ctk.CTkLabel(master=root,image=ctk.CTkImage(Image.open("./data/icon.png"),size=(85,85)),text="").place(x=5,y=5)
  tk.Label(root, bg="#0f1314", text="Midnight", fg="#fff", font=("Dubai Medium", 20)).place(x=100,y=0)
  tk.Label(root, bg="#0f1314", text="Raider", fg="#fff", font=("Dubai Medium", 20)).place(x=100,y=40)
  
  modulelist = ctk.CTkFrame(master=root, width=230, height=720, corner_radius=0, fg_color="#0f1314")
  modulelist.place(x=0,y=100)
  tk.Canvas(bg="#0a111e", highlightthickness=0, height=2080, width=4).place(x=230, y=0)
  ctk.CTkButton(master=modulelist, image=ctk.CTkImage(Image.open("data/join_leave.png"),size=(20, 20)), compound="left", fg_color="#0f1314", hover_color="#ade3f7", corner_radius=0, text="Joiner / Leaver", width=195, height=40, font=("Roboto", 16, "bold"), anchor="w", command= lambda: module_scroll_frame(1, 1)).place(x=20,y=12)
  #ctk.CTkButton(master=modulelist, image=ctk.CTkImage(Image.open("data/spammer.png"),size=(20, 20)), compound="left", fg_color="#0f1314", hover_color="#ade3f7", corner_radius=0, text=lang_load_set("spammer"), width=195, height=40, font=("Roboto", 16, "bold"), anchor="w", command= lambda: module_scroll_frame(1, 2)).place(x=20,y=57)
  ctk.CTkButton(master=modulelist, image=ctk.CTkImage(Image.open("data/setting.png"),size=(20, 20)), compound="left", fg_color="#0f1314", hover_color="#ade3f7", corner_radius=0, text="Settings", width=195, height=40, font=("Roboto", 16, "bold"), anchor="w", command= lambda: module_scroll_frame(2, 1)).place(x=20,y=516)
  ctk.CTkButton(master=modulelist, image=ctk.CTkImage(Image.open("data/info.png"),size=(20, 20)), compound="left", fg_color="#0f1314", hover_color="#ade3f7", corner_radius=0, text="About", width=195, height=40, font=("Roboto", 16, "bold"), anchor="w", command= lambda: module_scroll_frame(2, 2)).place(x=20,y=562)
  
  credit_frame = ctk.CTkFrame(root, width=1020, height=50, fg_color="#0a111c", bg_color="#152945")
  credit_frame.place(x=245, y=10)
  ctk.CTkButton(master=credit_frame, image=ctk.CTkImage(Image.open("data/link.png"),size=(20, 20)), compound="right", fg_color="#0a111c", text_color="#fff", corner_radius=0, text="", width=20, height=20, font=("Roboto", 16), anchor="w", command= lambda: CTkMessagebox(title="Version Info", message=f"Version: {version}\n\nDeveloper: {developer}\nTester: {testers}", width=450)).place(x=10,y=10)
  ctk.CTkLabel(master=credit_frame, fg_color="#0a111c", text_color="#fff", corner_radius=0, text="Username: "+os.getlogin(), width=20, height=20, font=("Roboto", 16, "bold"), anchor="w").place(x=40,y=5)
  ctk.CTkLabel(master=credit_frame, fg_color="#0a111c", text_color="#fff", corner_radius=0, text="Hwid: "+get_hwid(), width=20, height=20, font=("Roboto", 16, "bold"), anchor="w").place(x=40,y=25)
  ctk.CTkLabel(master=credit_frame, fg_color="#0a111c", text_color="#fff", corner_radius=0, text="Plan: "+get_plan(), width=20, height=20, font=("Roboto", 16, "bold"), anchor="w").place(x=450,y=25)


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
rpc_thread = threading.Thread(target=rpc.start)

0
root.mainloop()