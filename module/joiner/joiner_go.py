import subprocess
import re
import random
import threading

mentions = 0
allping = "None"

def start(token_file, proxie_file, serverid, invitelink, memberscreen, delay, bypasscaptcha, answers, apikey, deletejoinmsg, joinchannelid, useproxy, module_status):
    global process
    users = ['None']
    print("Starting the process.")
    print(delay)
    command = ['go', 'run', 'joiner_go.go', f'{token_file}', serverid, invitelink, str(memberscreen), f'{delay}', str(bypasscaptcha), answers, apikey, str(deletejoinmsg), joinchannelid, str(useproxy), f'{proxie_file}']
    #go run joiner_go.go C:/Users/Shin/Desktop/Data/GitHub/ThreeCoinRaider/module/spam/token_sample.txt 1197528360776126526 r33We25Z False 3 False None None False None True C:/Users/Shin/Desktop/Data/GitHub/ThreeCoinRaider/module/spam/proxie-lol.txt
    #// token_file serverid invitelink memberscreen delay bypasscaptcha answers apikey deletejoinmsg joinchannelid useproxy proxy_file
    print(command)
    process = subprocess.Popen(command, stdout=subprocess.PIPE, text=True, cwd=r"./module/spam/")
    monitor_thread = threading.Thread(target=monitor_process, args=(module_status, invitelink))
    monitor_thread.start()

def stop():
    global process
    print("Stopping the process.")
    if process.poll() is None:
        process.terminate()

def monitor_process(module_status, invitelink):
    global process
    while process.poll() is None:
        output = process.stdout.readline().strip()
        print(output)
        if output:
            matches = re.findall(r'\b\d+\b', output)
            if matches:
                channelid = matches[0]
            if '200' in output:
                print(f"[+] 参加に成功しました Invite: {invitelink}")
                module_status(1, 1, 1)
            if '400' in output:
                print(f"[-] 参加に失敗しました Invite: {invitelink}")
                module_status(1, 1, 2)