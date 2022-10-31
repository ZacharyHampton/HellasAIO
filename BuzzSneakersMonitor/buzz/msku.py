import time
import requests
import threading
from db import db
from buzz import webhook
from buzz.types import Thread
import sentry_sdk
import os
from dotenv import load_dotenv
from sentry_sdk import capture_exception

load_dotenv()
sentry_sdk.init(dsn=os.getenv('SENTRY_DSN'))


def BuzzMSKUFlow(MSKU: str, parentThread: Thread):
    print('[{}] Starting thread.'.format(MSKU))

    headers = {
        'Accept': 'application/json, text/javascript, */*; q=0.01',
        'Accept-Language': 'en-US,en;q=0.9',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive',
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
        'Origin': 'https://www.buzzsneakers.gr',
        'Pragma': 'no-cache',
        'Referer': 'https://www.buzzsneakers.gr/athlitika-papoutsia/',
        'Sec-Fetch-Dest': 'empty',
        'Sec-Fetch-Mode': 'cors',
        'Sec-Fetch-Site': 'same-origin',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36',
        'X-Requested-With': 'XMLHttpRequest',
        'sec-ch-ua': '".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"',
        'sec-ch-ua-mobile': '?0',
        'sec-ch-ua-platform': '"Windows"'
    }

    while not parentThread.stop:
        try:
            time.sleep(1)
            link = 'https://www.buzzsneakers.gr/files/images/slike_proizvoda/media/{}/{}/images/{}.jpg'.format(MSKU[:3], MSKU, MSKU)

            try:
                response = requests.get(
                    link,
                    headers=headers,
                )
            except requests.exceptions.ConnectionError:
                print('[{}] Connection error.'.format(MSKU))
                continue

            if response.status_code == 404:
                print('[{}] Product has not been loaded yet.'.format(MSKU))
                continue
            elif response.status_code == 200:
                print('[{}] Product has been loaded.'.format(MSKU))
                webhook.PingImageMSKU(MSKU, link)
                cur = db.cursor()
                cur.execute("DELETE FROM mskus WHERE msku = ?", (MSKU,))
                db.commit()
                cur.close()
                parentThread.Stop()
            else:
                print('[{}] Error: {}'.format(MSKU, response.status_code))
        except Exception as e:
            capture_exception(e)
            continue






