import subprocess
import re
import random
import threading

import module.spammer.utilities.user_scrape as user_scrape

mentions = 0
allping = "None"

def start(token_file, proxie_file, delay, tokens, module_status, serverid, channelid, contents, allchannel, allping, mentions, threads):
    global process
    users = ['None']
    print("Starting the process.")
    print(threads)
    if allping == True:
        users = user_scrape.get_members(serverid, channelid, random.choice(tokens))
        if users == None:
            print("[-] んーメンバーが取得できなかったっぽい token死なないように一回止めるね")
            return
        else:
            print(users)
        allping = "True"
    print(delay)
    command = ['go', 'run', 'spammer_go.go', serverid, channelid, contents, f'{token_file}', f'{proxie_file}', f'{threads}', f'{allchannel}', f'{delay}', f'{allping}', f'{int(mentions)}'] + users
    print(command)
    process = subprocess.Popen(command, stdout=subprocess.PIPE, text=True, cwd=r"./module/spammer/")
    monitor_thread = threading.Thread(target=monitor_process, args=(module_status, channelid))
    monitor_thread.start()

def stop():
    global process
    print("Stopping the process.")
    if process.poll() is None:
        process.terminate()

def monitor_process(module_status, channelid):
    global process
    while process.poll() is None:
        output = process.stdout.readline().strip()
        #print(output)
        if output:
            matches = re.findall(r'\b\d+\b', output)
            if matches:
                channelid = matches[0]
            if '[+]' in output:
                print(f"[+] 送信に成功しました ChannelID: {channelid}")
                module_status(2, 1, 1)
            elif '[-]' in output:
                print(f"[-] 送信に失敗しました ChannelID: {channelid}")
                module_status(2, 1, 2)
            else:
                print("Error:",output)