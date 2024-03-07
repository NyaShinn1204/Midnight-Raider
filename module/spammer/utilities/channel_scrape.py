import requests
import json
import utilities.header as header

def get_channels(token,guildid):
    global channels
    channels = []
    while True:
        req_header = header.request_header(token)
        headers = req_header
        x = requests.get(
            f"https://discord.com/api/v9/guilds/{guildid}/channels", headers=headers)
        data = json.loads(x.text)
        print(x.json)
        print(json.loads(x.text))
        if x.status_code == 200:
            print(x.status_code)
            for channel in data:
                if channel['type'] == 0 or 2:
                    if channel not in channels:
                        channels.append(channel["id"])
            return channels
        else:
            print(token)
            print(str(x.status_code))
            return