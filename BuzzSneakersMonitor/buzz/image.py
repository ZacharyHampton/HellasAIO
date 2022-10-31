import time
import requests
from buzz import webhook
from buzz.types import Thread
import sentry_sdk
import os
from dotenv import load_dotenv
from sentry_sdk import capture_exception
import xmltodict

load_dotenv()
sentry_sdk.init(dsn=os.getenv('SENTRY_DSN'))


def BuzzImageFlow(msku_partial: str, parentThread: Thread):
    print('[{}] Starting thread.'.format(msku_partial))

    headers = {
        'Accept': 'application/json, text/javascript, */*; q=0.01',
        'Accept-Language': 'en-US,en;q=0.9',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive',
        'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
        'Origin': 'https://www.buzzsneakers.gr',
        'Pragma': 'no-cache',
        'Referer': 'https://www.buzzsneakers.gr/nb-admin/lib/ckfinder/ckfinder.html',
        'Sec-Fetch-Dest': 'empty',
        'Sec-Fetch-Mode': 'cors',
        'Sec-Fetch-Site': 'same-origin',
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36',
        'X-Requested-With': 'XMLHttpRequest',
        'sec-ch-ua': '".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"',
        'sec-ch-ua-mobile': '?0',
        'sec-ch-ua-platform': '"Windows"'
    }

    currentMSKUS = []

    while not parentThread.stop:
        try:
            time.sleep(1)
            params = {
                'command': 'GetFolders',
                'type': 'Images',
                'currentFolder': '/slike_proizvoda/media/{}/'.format(msku_partial[:3]),
                'langCode': 'en',
                'hash': '4c2e6444e41bf385',
            }

            try:
                response = requests.get('https://www.buzzsneakers.gr/nb-admin/lib/ckfinder/core/connector/php/connector.php', params=params, headers=headers)
            except requests.exceptions.ConnectionError:
                print('[{}] Connection error.'.format(msku_partial))
                continue

            if response.status_code != 200:
                print('[{}] Error loading current loaded images.'.format(msku_partial))
                continue

            print('[{}] Got current images.'.format(msku_partial))

            data = xmltodict.parse(response.text)
            if parentThread.firstRun:
                for image in data['Connector']['Folders']['Folder']:
                    print('[{}] Adding image {}. (first run)'.format(msku_partial, image['@name']))
                    currentMSKUS.append(image['@name'])

                parentThread.firstRun = False
                continue

            for image in data['Connector']['Folders']['Folder']:
                msku = image['@name']

                if msku_partial not in msku:
                    # print('[{}] Skipping image {}. It does not match the format.'.format(msku_partial, msku))
                    continue

                if msku not in currentMSKUS:
                    print('[{}] New image found: {}'.format(msku_partial, msku))
                    imageUrl = 'https://www.buzzsneakers.gr/files/images/slike_proizvoda/media/{}/{}/images/{}.jpg'.format(
                        msku_partial, msku, msku)
                    currentMSKUS.append(msku)
                    webhook.PingImageMSKU(msku, imageUrl)
                else:
                    print('[{}] Image already found: {}'.format(msku_partial, msku))
                    continue
        except Exception as e:
            capture_exception(e)
            continue
