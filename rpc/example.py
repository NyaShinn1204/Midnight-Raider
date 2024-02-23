import rpc
import time
from time import mktime

print("RPC connection successful.")
time.sleep(5)
start_time = mktime(time.localtime())
while True:
    activity = {
            "state": "Midnight Raider v1.0.0",  # anything you like
            "details": "Raiding a Server",  # anything you like
            "timestamps": {
                "start": start_time
            },
            "assets": {
                #"small_text": "Midnight",  # anything you like
                #"small_image": "icon",  # must match the image key
                "large_text": "Midnight Raider",  # anything you like
                "large_image": "icon"  # must match the image key
            }
        }
    rpc.DiscordIpcClient.for_platform("1210561968969220218").set_activity(activity)
    time.sleep(900)