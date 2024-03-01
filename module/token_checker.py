import threading
from httpx import Client

def check(tokens, update_token):
    lock = threading.Lock()
    def success(text):
        lock.acquire()
        update_token(True, text)
        lock.release()
    def invalid(text):
        lock.acquire()
        update_token(False, text)
        lock.release()
    def check_token(token:str):
        request = Client()
        x = request.get('https://discord.com/api/v9/users/@me/library', headers={"authorization": token},timeout=5)
        if x.status_code == 200:
            success(token)
        else:
            invalid(token)
    def check_tokens():
        threads=[]
        for token in tokens:
            try:
                threads.append(threading.Thread(target=check_token, args=(token,)))
            except Exception as e:
                pass
        for thread in threads:
             thread.start()
        for thread in threads:
             thread.join()
    threading.Thread(target=check_tokens).start()