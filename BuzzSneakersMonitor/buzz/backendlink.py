import time
import requests
from buzz.utils import GetPIDFromLink, GetStockBoolSize
from buzz import webhook
from buzz.types import Thread, Size
import sentry_sdk
import os
from dotenv import load_dotenv
import xmltodict
from sentry_sdk import capture_exception

load_dotenv()
sentry_sdk.init(dsn=os.getenv('SENTRY_DSN'))


def BackendLinkFlow(_, parentThread: Thread):
    print('[BACKENDLINK] Starting thread.')

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
    links = []

    while not parentThread.stop:
        try:
            time.sleep(1)
            try:
                response = requests.get(
                    'https://www.buzzsneakers.gr/files/sitemap/GRC_gr/product.xml',
                    headers=headers,
                )
            except requests.exceptions.ConnectionError or requests.exceptions.ConnectTimeout:
                print("[BACKENDLINK] Connection error.")
                continue

            if response.status_code == 200:
                try:
                    data = xmltodict.parse(response.text)
                except xmltodict.expat.ExpatError:
                    print('[BACKENDLINK] XML Error.')
                    continue

                print("[BACKENDLINK] Successfully fetched data.")

                if not data['urlset'].get('url'):
                    print("[BACKENDLINK] No links found.")
                    continue

                if parentThread.firstRun:
                    parentThread.firstRun = False

                    for url in data['urlset']['url']:
                        links.append(url['loc'])
                        print("[BACKENDLINK] Added link: {}".format(url['loc']))

                    continue

                if [url['loc'] for url in data['urlset']['url']] == links:
                    print("[BACKENDLINK] No new links found.")
                    continue

                for url in data['urlset']['url']:
                    if url['loc'] not in links:
                        links.append(url['loc'])
                        print("[BACKENDLINK] Added link: {}".format(url['loc']))
                        pid = GetPIDFromLink(url['loc'])

                        while True:
                            print("[BACKENDLINK - {}] Getting product info.".format(pid))
                            sizes = []
                            try:
                                response = requests.post(
                                    'https://www.buzzsneakers.gr/athlitika-papoutsia/',
                                    headers=headers,
                                    data="nbAjax=1&task=getproductdata&productId={}".format(pid)
                                )
                            except requests.exceptions.ConnectionError:
                                print("[BACKENDLINK - {}] Connection error.".format(pid))
                                continue


                            if response.status_code != 200:
                                time.sleep(1)
                                continue

                            rJSON = response.json()
                            for size in rJSON['sizes']:
                                sizes.append(Size(
                                    pMSKU=rJSON['product']['productCode'],
                                    stock=int(float(size['quantity'])),
                                    size=size['sizeName'],
                                    instock=bool(GetStockBoolSize(size))
                                ))

                            webhook.PingServer(
                                webhookTitle="Product stock change detected!",
                                productName=rJSON['product']['name'],
                                productMSKU=rJSON['product']['productCode'],
                                productSizes=sizes,
                                productImage="https://www.buzzsneakers.gr/{}".format(rJSON['product']['image']),
                                sizeTitle="Sizes",
                                productId=pid,
                                productLink=rJSON['product']['permalinkProduct'],
                                includePrice=False,
                                quantity=int(float(rJSON['product']['quantity'])),
                            )

                            break

                links = [url['loc'] for url in data['urlset']['url']]
            else:
                print("[BACKENDLINK] Error fetching data.")
        except Exception as e:
            print('[BACKENDLINK] Error: {}'.format(e))
            capture_exception(e)
            continue
