import random
import string
import threading
import time
import json
import re
from httpx import Client
from httpx_socks import SyncProxyTransport

import module.spammer.utilities.channel_scrape as ch_scrape
import module.spammer.utilities.user_scrape as user_scrape
import utilities.header as header
import utilities.random_convert as random_convert

status = True
timelock = False

def status():
    global status
    return status

def stop():
    global status
    status = False

def randomname(n):
    return ''.join(random.choice(string.ascii_letters + string.digits) for i in range(n))  

def extract(format_token):
    if re.compile(r"(.+):").match(format_token):
        return format_token.split(":")[1]
    else:
        token = format_token
    return token

def start(delay, tokens, module_status, proxysetting, proxies, proxytype, serverid, channelid, contents, allchannel, allping, mentions, randomstring, ratelimit, randomconvert):
    global status
    global channels
    global users
    global timelock
    status = True
    
    token = random.choice(tokens)
    
    #print(token)
    #print(serverid)
    
    # Get Channel, User
    if allchannel == True:
        channels = ch_scrape.get_channels(token,int(serverid))
        if channels == None:
            print("[-] んーチャンネルが取得できなかったっぽい token死なないように一回止めるね")
            return
    if allping == True:
        users = user_scrape.get_members(serverid, channelid, token)
        if users == None:
            print("[-] んーメンバーが取得できなかったっぽい token死なないように一回止めるね")
            return
    
    while status is True:
        if status == False:
            break
        if timelock == True:
            print("[-] RateLimit Fixing...")
            time.sleep(8)
            print("[+] RateLimit Fixed")
            timelock = False
        threading.Thread(target=spammer_thread, args=(tokens, module_status, proxysetting, proxies, proxytype,
                        allchannel, allping, channelid, contents, randomstring, mentions, ratelimit, randomconvert)).start()
        time.sleep(float(delay))
        
def spammer_thread(tokens, module_status, proxysetting, proxies, proxytype, allchannel, allping, channelid, contents, randomstring, mentions, ratelimit, randomconvert):
    global channels
    global users
    global status
    global timelock

    if timelock == True:
        return
    if status is False:
        return
    token = random.choice(tokens)
    content = contents
    if content == "":
        print("[-] メッセージが設定されていないので初期のメッセージを送信します")
        content = "ThreeCoinRaider On Top :skull:"
    if randomconvert == True:
        content = random_convert.random_convert(content)
    if allping == True:
        for i in range(int(mentions)):
            content = content + f"<@{random.choice(users)}>"
    if randomstring == True:
        content = f"{content}\n{randomname(10)}"
    if allchannel == True:
        channelid = random.choice(channels)
    data = {"content": content}
    req_header = header.request_header(token)
    headers = req_header
    extract_token = f"{extract(token+']').split('.')[0]}.{extract(token+']').split('.')[1]}"
    try:
        if status is False:
            return
        request = Client()
        if proxysetting == True:
            proxy = random.choice(proxies)
            request = Client(transport=SyncProxyTransport.from_url(f'{proxytype}://{proxy}'))
        x = request.post(f"https://discord.com/api/v9/channels/{channelid}/messages", headers=headers, json=data)
        if x.status_code == 400:
            print(f"[-] AutoModによりメッセージが削除されたっぽい  Message: {x.json()['message']} ChannelID: {channelid} Token: {extract_token}.******** Status: {x.status_code}")
            module_status(2, 1, 2)
        if x.status_code == 403:
            print(f"[-] このチャンネルで発現できません ChannelID: {channelid} Token: {extract_token}.******** Status: {x.status_code}")
            module_status(2, 1, 2)
        if x.status_code == 404:
            print(f"[-] このチャンネルは存在しません ChannelID: {channelid} Token: {extract_token}.******** Status: {x.status_code}")
            module_status(2, 1, 2)
        if x.status_code == 200:
            module_status(2, 1, 1)
            if proxysetting == True:
                print(f"[+] 送信に成功しました ChannelID: {channelid} Token: {extract_token}.******** Proxy: {proxy}")
            else:
                print(f"[+] 送信に成功しました ChannelID: {channelid} Token: {extract_token}.********")
        else:
            if x.status_code == 429 or x.status_code == 20016:
                print("[-] RateLimit!! Please Wait!! "+x.json()["retry_after"])
                if ratelimit == True:
                    timelock = True
                return
            module_status(2, 1, 2)
    except:
        pass