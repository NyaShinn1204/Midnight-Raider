import rpc as rpc
import time 
from time import mktime

def start():
    client_id = '1210561968969220218'  # Your application's client ID as a string. (This isn't a real client ID)
    rpc_obj = rpc.DiscordIpcClient.for_platform(client_id)  # Send the client ID to the rpc module
    #print("RPC connection successful.")
    
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
        rpc_obj.set_activity(activity)
        time.sleep(900)
        
start()