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

# Module Import
import module.joiner.joiner as module_joiner

# Utilities Module Import
import module.token_checker as token_checker
import module.proxy_checker as proxy_checker

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
    
def extractfi(input_str):
  if len(input_str) >= 5:
    replaced_str = input_str[:-5] + '*' * 5
    return replaced_str
  else:
    return input_str
  
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

def get_invite(invite):
  try:
    x = requests.get(f"https://discord.com/api/v9/invites/{invite}?with_counts=true&with_expiration=true")
    if x.status_code == 200:
      verify_level = x.json()["guild"]["vverification_level"]
      if verify_level == 1:
        printl("info", "Only Email Verify members")
      if verify_level == 2:
        printl("info", "Only members who have been on the server for 5 minutes")
      if verify_level == 3:
        printl("info", "Only members who have been on the server for 10 minutes")
      if verify_level == 4:
        printl("info", "Only Phone Verify members")
    return x.status_code
  except:
    printl("error", "get_invite error wrong")

def gui_close():
  root.destroy()
  rpc.stop_threads = True

System.Size(120, 30)
System.Clear()

# c2, c5, c7, c4

#c1 = "#0D2845"
c2 = "#020b1f"
c3 = "#b4e7fa"
#c4 = "#020b1f"
c5 = "#141877"
c6 = "#b4e7fa"
c7 = "#000117"
#c8 = "#489ea1"
#c9 = "#454c7f"
#c10 = "#2D2DA0"
#c11 = "#041432"
#c12 = "#3a88e3"
c13 = "#141B39"

root = tk.Tk()
root.geometry("1280x720")
root.resizable(0, 0)
root.title("Midnight Raider | "+version)
root.iconbitmap("./data/icon.ico")
root.configure(bg="#baebfb")
root.protocol("WM_DELETE_WINDOW", gui_close)

# Import Variable
from data.settings import Setting, SettingVariable

# Import Utilities
import module.joiner.utilities.get_balance as get_balance

def load_background():
  ctk.CTkLabel(master=root,image=ctk.CTkImage(Image.open("./data/background-01.jpg"),size=(1280,720)),text="").pack()

# Check Config
def check_config():
  printl("info", "Checking Config")
  try:
    if os.path.exists(r"config.json") and not json.load(open('./config.json', 'r', encoding="utf-8"))["token_path"] == "":
      tokens = open(json.load(open('./config.json', 'r', encoding="utf-8"))["token_path"], 'r').read().splitlines()
      Setting.tokens = []
      Setting.validtoken = 0
      Setting.invalidtoken = 0
      Setting.token_filenameLabel.set(os.path.basename(json.load(open('./config.json', 'r', encoding="utf-8"))["token_path"]))
      Setting.totaltokenLabel.set("Total: "+str(len(tokens)).zfill(3))
      threading.Thread(target=token_checker.check(tokens, update_token)).start()
      printl("info", "Checked Config")
    else:
      printl("error", "Config Not Found")
      printl("error", "Please point to it manually.")
      token_load()
  except Exception as error:
    printl("error", "Config Check Error")
    printl("error", error)
    token_load()

# Load Token, Proxie
# Token Tab
def token_load():
  filepath = filedialog.askopenfilename(filetype=[("", "*.txt")], initialdir=os.path.abspath(os.path.dirname(__file__)), title="Select Tokens")
  if filepath == "":
    return
  tokens = open(filepath, 'r').read().splitlines()
  if tokens == []:
    return
  data = json.load(open('config.json'))
  data['token_path'] = filepath
  json.dump(data, open('config.json', 'w'), indent=4)
  printl("info", f"Set Token File {os.path.basename(filepath)}")
  Setting.tokens = []
  Setting.validtoken = 0
  Setting.invalidtoken = 0
  Setting.token_filenameLabel.set(os.path.basename(filepath))
  Setting.validtokenLabel.set("Valid: 000")
  Setting.invalidtokenLabel.set("Invalid: 000")
  Setting.totaltokenLabel.set("Total: "+str(len(tokens)).zfill(3))
  threading.Thread(target=token_checker.check(tokens, update_token)).start()

def update_token(status, token):
  if status == True:
    Setting.tokens.append(token)
    Setting.validtoken += 1
    Setting.validtokenLabel.set("Valid: "+str(Setting.validtoken).zfill(3))
  if status == False:
    Setting.invalidtoken += 1
    Setting.invalidtokenLabel.set("Invalid: "+str(Setting.invalidtoken).zfill(3))

# Proxy Tab
def proxy_load():
  threading.Thread(target=proxy_main).start()
  
def proxy_main():
  proxy_type = Setting.proxytype.get()
  print(proxy_type)
  if proxy_type == "":
    print("[-] Cancel proxy")
    return
  proxy_filepath()

def proxy_filepath():
  filepath = filedialog.askopenfilename(filetype=[("", "*.txt")], initialdir=os.path.abspath(os.path.dirname(__file__)), title="Select Proxies")
  if filepath == "":
    return
  proxies = open(filepath, 'r').read().splitlines()
  if proxies == []:
    return
  data = json.load(open('config.json'))
  data['proxie_path'] = filepath
  json.dump(data, open('config.json', 'w'), indent=4)
  printl("info", f"Set Proxie File {os.path.basename(filepath)}")
  Setting.proxies = []
  Setting.totalproxies = str(len(proxies))
  Setting.vaildproxies = 0
  Setting.invaildproxies = 0
  Setting.proxy_filenameLabel.set(os.path.basename(filepath))
  Setting.totalProxiesLabel.set("Total: "+Setting.totalproxies.zfill(3))
  print("[+] Load: " + Setting.totalproxies + "Proxies")
  time.sleep(1)
  threading.Thread(target=proxy_checker.check(update_proxy, proxies, Setting.proxytype.get()))
  if Setting.vaildproxies == 0:
    printl("error","Not Found Load Vaild Proxies")
  else:
    printl("info","Success Load Vaild Proxies: " + str(Setting.vaildproxies))
     
def update_proxy(status, proxy):
  if status == True:
    Setting.proxies.append(proxy)
    Setting.vaildproxies += 1
    Setting.validProxiesLabel.set("Valid: "+str(Setting.vaildproxies).zfill(3))
  if status == False:
    Setting.invaildproxies += 1
    Setting.invalidProxiesLabel.set("Invalid: "+str(Setting.invaildproxies).zfill(3))

def clear_frame(frame):
  for widget in frame.winfo_children():
    widget.destroy()
  frame.pack_forget()

def module_thread(num1, num2, num3):
  tokens = Setting.tokens
  proxies = Setting.proxies
  proxytype = Setting.proxytype.get()
  proxysetting = Setting.proxy_enabled.get()
  delay = 0.1
  
  if num1 == 1:
    if num2 == 1:
      if num3 == 1:
        serverid = str(Setting.joiner_serverid.get())
        join_channelid = str(Setting.joiner_channelid.get())
        invitelink = str(Setting.joiner_link.get())
        memberscreen = Setting.joiner_bypassms.get()
        delete_joinms = Setting.joiner_deletems.get()
        bypasscaptcha = Setting.joiner_bypasscap.get()
    
        delay = Setting.joiner_delay.get()
    
        answers = None
        api = None
    
        if invitelink == "":
          print("[-] InviteLink is not set")
          return
        if invitelink.__contains__('discord.gg/'):
          invitelink = invitelink.replace('discord.gg/', '').replace('https://', '').replace('http://', '')
        elif invitelink.__contains__('discord.com/invite/'):
          invitelink = invitelink.replace('discord.com/invite/', '').replace('https://', '').replace('http://', '')
        try:
          invitelink = invitelink.split(".gg/")[1]
        except:
          pass
        if memberscreen == True:
          if serverid == "":
            print("[-] ServerID is not set")
          else:
            print("[-] このオプションは非推奨です")
        if bypasscaptcha == True:
          if answers == "":
            print("[-] Please Select API Service")
            return
          else:
            if api == "":
              print("[-] Please Input API Keys")
        if delete_joinms == True:
          if join_channelid == "":
            print("[-] Join ChannelID is not set")
            return
        
        if get_invite(invitelink) == 404:
          printl("error", "This invite code not found")
          return  
        
        threading.Thread(target=module_joiner.start, args=(tokens, serverid, invitelink, memberscreen, delay, module_status, answers, api, bypasscaptcha, delete_joinms, join_channelid)).start()

def module_status(num1, num2, num3):
  if num1 == 1:
    if num2 == 1:
      if num3 == 1:
        SettingVariable.joinerresult_success +=1
        Setting.suc_joiner_Label.set("Success: "+str(SettingVariable.joinerresult_success).zfill(3))
      if num3 == 2:
        SettingVariable.joinerresult_failed +=1
        Setting.fai_joiner_Label.set("Failed: "+str(SettingVariable.joinerresult_failed).zfill(3))
  if num1 == 2:
    if num2 == 1:
      if num3 == 1:
        print("2-1-1")
      if num3 == 2:
        print("2-1-2")

def module_scroll_frame(num1, num2):
  global module_frame
  frame_scroll = module_frame = ctk.CTkScrollableFrame(root, fg_color="#0f1314", bg_color="#152945", width=1000, height=630)
  module_frame.place(x=245, y=70)
  clear_frame(frame_scroll)
  if num1 == 1:
    if num2 == 1:
      # Joiner
      # Frame Number 01_01
      def hcaptcha_select():
        global answers, api
        if Setting.joiner_bypasscap.get() == True:
          answers = ctk.CTkInputDialog(text = "Select Sovler\n1, CapSolver\n2, CapMonster\n3, 2Cap\n4, Anti-Captcha").get_input()
          if answers in ['1','2','3','4']:
            print("[+] Select " + answers)
            api = ctk.CTkInputDialog(text = "Input API Key").get_input()
            if api == "":
              print("[-] Not Set. Please Input")
              Setting.joiner_bypasscap.set(False)
            else:
              print("[~] Checking API Key: " + extractfi(api))
              if answers == "1":
                if get_balance.get_balance_capsolver(api) == 0.0:
                  Setting.joiner_bypasscap.set(False)
              if answers == "2":
                if get_balance.get_balance_capmonster(api) == 0.0:
                  Setting.joiner_bypasscap.set(False)
              if answers == "3":
                if get_balance.get_balance_2cap(api) == 0.0:
                  Setting.joiner_bypasscap.set(False)
              if answers == "4":
                if get_balance.get_balance_anticaptcha(api) == 0.0:
                  Setting.joiner_bypasscap.set(False)
          else:
            print("[-] Not Set. Please Input")
            Setting.joiner_bypasscap.set(False)

      modules_frame01_01 = ctk.CTkFrame(module_frame, width=470, height=275, border_width=0, fg_color=c13)
      modules_frame01_01.grid(row=0, column=0, padx=6, pady=6)
      tk.Label(modules_frame01_01, bg=c13, fg="#fff", text="Joiner", font=("Roboto", 12, "bold")).place(x=15,y=0)
      tk.Canvas(modules_frame01_01, bg=c6, highlightthickness=0, height=4, width=470).place(x=0, y=25)
      
      ctk.CTkCheckBox(modules_frame01_01, bg_color=c13, text_color="#fff", border_color=c3, checkbox_width=20, checkbox_height=20, hover=False, border_width=3, text="Bypass MemberScreen", variable=Setting.joiner_bypassms).place(x=5,y=31)
      test = ctk.CTkLabel(modules_frame01_01, text_color="#fff", text="(?)")
      test.place(x=170,y=31)
      CTkToolTip(test, delay=0.5, message="Bypass the member screen when you join.") 
      ctk.CTkCheckBox(modules_frame01_01, bg_color=c13, text_color="#fff", border_color=c3, checkbox_width=20, checkbox_height=20, hover=False, border_width=3, text="Bypass hCaptcha", variable=Setting.joiner_bypasscap, command=hcaptcha_select).place(x=5,y=55) 
      test = ctk.CTkLabel(modules_frame01_01, text_color="#fff", text="(?)")
      test.place(x=140,y=55)
      CTkToolTip(test, delay=0.5, message="Automatically resolve hcaptcha")
      ctk.CTkCheckBox(modules_frame01_01, bg_color=c13, text_color="#fff", border_color=c3, checkbox_width=20, checkbox_height=20, hover=False, border_width=3, text="Delete Join Message", variable=Setting.joiner_deletems).place(x=5,y=79)
      test = ctk.CTkLabel(modules_frame01_01, text_color="#fff", text="(?)")
      test.place(x=160,y=79)
      CTkToolTip(test, delay=0.5, message="Delete the message when you join") 
      
      ctk.CTkButton(modules_frame01_01, text="Clear        ", fg_color=c2, hover_color=c5, width=75, height=25, command=lambda: Setting.joiner_link.set("")).place(x=5,y=109)
      ctk.CTkEntry(modules_frame01_01, bg_color=c13, fg_color=c7, border_color=c2, text_color="#fff", width=150, height=20, textvariable=Setting.joiner_link).place(x=85,y=109)
      tk.Label(modules_frame01_01, bg=c13, fg="#fff", text="Invite Link", font=("Roboto", 12)).place(x=240,y=107)
      ctk.CTkButton(modules_frame01_01, text="Clear        ", fg_color=c2, hover_color=c5, width=75, height=25, command=lambda: Setting.joiner_serverid.set("")).place(x=5,y=138)
      ctk.CTkEntry(modules_frame01_01, bg_color=c13, fg_color=c7, border_color=c2, text_color="#fff", width=150, height=20, textvariable=Setting.joiner_serverid).place(x=85,y=138)
      tk.Label(modules_frame01_01, bg=c13, fg="#fff", text="Server ID", font=("Roboto", 12)).place(x=240,y=136)
      test = ctk.CTkLabel(modules_frame01_01, text_color="#fff", text="(?)")
      test.place(x=320,y=136)
      CTkToolTip(test, delay=0.5, message="Used on the member screen, \nbut locked before the member screen is bypassed") 
      ctk.CTkButton(modules_frame01_01, text="Clear        ", fg_color=c2, hover_color=c5, width=75, height=25, command=lambda: Setting.joiner_channelid.set("")).place(x=5,y=167)
      ctk.CTkEntry(modules_frame01_01, bg_color=c13, fg_color=c7, border_color=c2, text_color="#fff", width=150, height=20, textvariable=Setting.joiner_channelid).place(x=85,y=167)
      tk.Label(modules_frame01_01, bg=c13, fg="#fff", text="Channel ID", font=("Roboto", 12)).place(x=240,y=165)

      CTkLabel(modules_frame01_01, text_color="#fff", text="Delay Time (s)", font=("Roboto", 15)).place(x=5,y=192)
      def show_value01_01_01(value):
          tooltip01_01_01.configure(message=round(value, 1))
      test = ctk.CTkSlider(modules_frame01_01, from_=0.1, to=3.0, variable=Setting.joiner_delay, command=show_value01_01_01)
      test.place(x=5,y=217)
      tooltip01_01_01 = CTkToolTip(test, message=round(Setting.joiner_delay.get(), 1))

      ctk.CTkButton(modules_frame01_01, text="Start", fg_color=c2, hover_color=c5, width=60, height=25, command=lambda: module_thread(1, 1, 1)).place(x=5,y=237)

      tk.Label(modules_frame01_01, bg=c13, fg="#fff", text="Join Status", font=("Roboto", 12)).place(x=205,y=190)
      tk.Label(modules_frame01_01, bg=c13, fg="#fff", textvariable=Setting.suc_joiner_Label, font=("Roboto", 12)).place(x=210,y=215)
      tk.Label(modules_frame01_01, bg=c13, fg="#fff", textvariable=Setting.fai_joiner_Label, font=("Roboto", 12)).place(x=210,y=240)
  
      printl("info", "Open Join Leave Tab")
  if num1 == 2:
    if num2 == 1:
      # Setting

      # Frame Numnber 10_01
      modules_frame10_01 = ctk.CTkFrame(module_frame, width=470, height=210, border_width=0, fg_color=c13)
      modules_frame10_01.grid(row=0, column=0, padx=6, pady=6)
      tk.Label(modules_frame10_01, bg=c13, fg="#fff", text="Tokens", font=("Roboto", 12, "bold")).place(x=15,y=0)
      tk.Canvas(modules_frame10_01, bg=c6, highlightthickness=0, height=4, width=470).place(x=0, y=25)

      ctk.CTkButton(modules_frame10_01, text="Select File", fg_color=c2, hover_color=c5, width=75, height=25, command=lambda: token_load(), font=("Roboto", 12)).place(x=5,y=33)
      ctk.CTkEntry(modules_frame10_01, bg_color=c13, fg_color=c7, border_color=c2, text_color="#fff", width=150, height=20, state="disabled").place(x=85,y=33)
      ctk.CTkLabel(modules_frame10_01, bg_color=c13, fg_color=c2, text_color="#fff", text="", width=150, height=20, textvariable=Setting.token_filenameLabel).place(x=85,y=33)
      tk.Label(modules_frame10_01, bg=c13, fg="#fff", text="File Name", font=("Roboto", 12)).place(x=240,y=31)

      tk.Label(modules_frame10_01, bg=c13, fg="#fff", text="Status", font=("Roboto", 12)).place(x=5,y=70)
      tk.Label(modules_frame10_01, bg=c13, fg="#fff", text="Total: 000", font=("Roboto", 12), textvariable=Setting.totaltokenLabel).place(x=10,y=95)
      tk.Label(modules_frame10_01, bg=c13, fg="#fff", text="Valid: 000", font=("Roboto", 12), textvariable=Setting.validtokenLabel).place(x=10,y=115)
      tk.Label(modules_frame10_01, bg=c13, fg="#fff", text="Invalid: 000", font=("Roboto", 12), textvariable=Setting.invalidtokenLabel).place(x=10,y=135)
      
      
      # Frame Numnber 10_02
      modules_frame10_02 = ctk.CTkFrame(module_frame, width=470, height=210, border_width=0, fg_color=c13)
      modules_frame10_02.grid(row=0, column=1, padx=6, pady=6)
      tk.Label(modules_frame10_02, bg=c13, fg="#fff", text="Proxies", font=("Roboto", 12, "bold")).place(x=15,y=0)
      tk.Canvas(modules_frame10_02, bg=c6, highlightthickness=0, height=4, width=470).place(x=0, y=25)

      ctk.CTkCheckBox(modules_frame10_02, bg_color=c13, text_color="#fff", border_color=c3, checkbox_width=20, checkbox_height=20, hover=False, border_width=3, variable=Setting.proxy_enabled, text="Enabled").place(x=5,y=31)
      def set_socket(socks):
        Setting.proxytype.set(socks)
      ctk.CTkOptionMenu(modules_frame10_02, height=25, corner_radius=4, values=["http", "https", "socks4", "socks5"], fg_color="#0a111c", button_color="#0a111c", button_hover_color="#0a111c", dropdown_fg_color="#0a111c", dropdown_hover_color=c13, dropdown_text_color="#fff", font=("Roboto", 12, "bold"), dropdown_font=("Roboto", 12, "bold"), command=set_socket, variable=Setting.proxytype).place(x=5,y=57)
      

      
      tk.Label(modules_frame10_02, bg=c13, fg="#fff", text="Socket Type", font=("Roboto", 12)).place(x=150,y=55)
      ctk.CTkButton(modules_frame10_02, text="Select File", fg_color=c2, hover_color=c5, width=75, height=25, command=lambda: proxy_load()).place(x=5,y=90)
      ctk.CTkEntry(modules_frame10_02, bg_color=c13, fg_color=c7, border_color=c2, text_color="#fff", width=150, height=20, state="disabled").place(x=85,y=90)
      ctk.CTkLabel(modules_frame10_02, bg_color=c13, fg_color=c2, text_color="#fff", text="", width=150, height=20, textvariable=Setting.proxy_filenameLabel).place(x=85,y=90)
      tk.Label(modules_frame10_02, bg=c13, fg="#fff", text="File Name", font=("Roboto", 12)).place(x=240,y=87)
    
      tk.Label(modules_frame10_02, bg=c13, fg="#fff", text="Status", font=("Roboto", 12)).place(x=5,y=120)
      tk.Label(modules_frame10_02, bg=c13, fg="#fff", text="Total: 000", font=("Roboto", 12), textvariable=Setting.totalProxiesLabel).place(x=10,y=145)
      tk.Label(modules_frame10_02, bg=c13, fg="#fff", text="Valid: 000", font=("Roboto", 12), textvariable=Setting.validProxiesLabel).place(x=10,y=165)
      tk.Label(modules_frame10_02, bg=c13, fg="#fff", text="Invalid: 000", font=("Roboto", 12), textvariable=Setting.invalidProxiesLabel).place(x=10,y=185)
   
      
      printl("info", "Open Setting Tab")
    if num2 == 2:
      # About
      credits_frame = ctk.CTkFrame(module_frame, width=940, height=400, border_width=0, corner_radius=5, fg_color="#0f1314")
      credits_frame.grid(row=1, column=1, padx=6, pady=6)
      tk.Label(credits_frame, bg="#0f1314", fg="#fff", text="Midnight Raider github:", font=("Roboto", 12)).place(x=0,y=0)
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

      printl("info", "Open About Tab")

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
  ctk.CTkButton(master=credit_frame, image=ctk.CTkImage(Image.open("data/link.png"),size=(20, 20)), compound="right", fg_color="#0a111c", hover_color="#0a111c", text_color="#fff", corner_radius=0, text="", width=20, height=20, font=("Roboto", 16), anchor="w", command= lambda: CTkMessagebox(title="Version Info", message=f"Version: {version}\n\nDeveloper: {developer}\nTester: {testers}", width=450)).place(x=10,y=10)
  ctk.CTkLabel(master=credit_frame, fg_color="#0a111c", text_color="#fff", corner_radius=0, text="Username: "+os.getlogin(), width=20, height=20, font=("Roboto", 16, "bold"), anchor="w").place(x=40,y=5)
  ctk.CTkLabel(master=credit_frame, fg_color="#0a111c", text_color="#fff", corner_radius=0, text="Hwid: "+get_hwid(), width=20, height=20, font=("Roboto", 16, "bold"), anchor="w").place(x=40,y=25)


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
check_config()
module_list_frame()
printl("info", "Starting RPC")
rpc_thread = threading.Thread(target=rpc.start).start()


root.mainloop()