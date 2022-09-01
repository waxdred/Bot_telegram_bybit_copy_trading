from telethon import TelegramClient
from dotenv import load_dotenv
from pathlib import Path

import os


dotenv_path = Path('./.env')
load_dotenv(dotenv_path=dotenv_path)
if os.path.exists("_log"):
    os.remove("_log")

session = "trading bot"
api_id = os.getenv('API_ID')
api_hash = os.getenv('API_HASH')
id_channel = os.getenv('ID_CHANNEL')
proxy = None

async def handler(update):
    t = str(update)
    pos = t.find(str(id_channel))
    t = t[pos:]
    pos = t.find("message=")
    if pos != -1:
        t = t[pos + len("message= "):]
        pos = t.find("'")
        if pos != -1:
            if os.path.exists("_log"):
                os.remove("_log")
            t = t[:pos]
            print(t)
            f = open("_log", "w")
            f.write(t)
            f.close()


# Use the client in a `with` block. It calls `start/disconnect` automatically.
with TelegramClient(session, api_id, api_hash, proxy=proxy) as client:
    # Register the update handler so that it gets called
    client.add_event_handler(handler)

    # Run the client until Ctrl+C is pressed, or the client disconnects
    client.run_until_disconnected()

f.close()
