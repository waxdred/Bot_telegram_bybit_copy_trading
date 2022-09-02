from telethon import TelegramClient
from dotenv import load_dotenv
from pathlib import Path

import os

from telethon.tl.custom import message

dotenv_path = Path('./.env')
load_dotenv(dotenv_path=dotenv_path)

class Telegram():
    def __init__(self):
        self.HOST = '127.0.0.1'
        self.PORT = 30000
        self.api_id = os.getenv('API_ID')
        self.api_hash = os.getenv('API_HASH')
        self.id_channel = os.getenv('ID_CHANNEL')
        self.my_channel = os.getenv('MY_CHANNEL')
        self.session = "trading bot"
        self.proxy = None
        self.msg = ""
        self.client = TelegramClient(self.session, self.api_id, self.api_hash, proxy=self.proxy)

    async def handler(self, update):
        t = str(update)
        pos = t.find(str(self.id_channel))
        t = t[pos:]
        pos = t.find("message=")
        if pos != -1:
            t = t[pos + len("message= "):]
            pos = t.find("'")
            if pos != -1:
                if os.path.exists("_log"):
                    os.remove("_log")
                t = t[:pos]
                t = t.replace("\\n", "\n", -1)
                print(t)
                await self.client.send_message(entity=int(self.my_channel), message=t)
        return t

    def start(self):
        with self.client:
            # Register the update handler so that it gets called
            t = self.client.add_event_handler(self.handler)
            print(t)
        
            # Run the client until Ctrl+C is pressed, or the client disconnects
            self.client.run_until_disconnected()

if __name__ == '__main__':
    print("Run Python")
    telegram = Telegram()
    print("Class init....")
    telegram.start()
