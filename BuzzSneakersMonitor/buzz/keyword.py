import time
import requests
from buzz.utils import GetPIDFromLink, GetStockBoolSize
from buzz import webhook
from buzz.types import Thread, Size
import sentry_sdk
import os
from dotenv import load_dotenv
from bs4 import BeautifulSoup
from sentry_sdk import capture_exception

load_dotenv()
sentry_sdk.init(dsn=os.getenv('SENTRY_DSN'))


def KeywordFlow(keywords: list, parentThread: Thread):
    print('[KEYWORDS] Starting thread.')

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
    currentProducts = []

    while not parentThread.stop:
        try:
            time.sleep(1)
            data = {
                'listId': '',
                'prices': '',
                'search': '',
                'sizeEU': '',
                'sort': 'new',
                'limit': '64',
                'typeView': 'grid',
                'typesearch': '1',
                'size': '',
                'sale': 'no',
                'outlet': 'no',
                'page': '0',
                'page_url': 'https://www.buzzsneakers.gr/athlitika-papoutsia/',
                'author': '',
                'ajax': 'yes',
                'separateAjax': 'undefined',
            }
            try:
                response = requests.post('https://www.buzzsneakers.gr/athlitika-papoutsia/', headers=headers, data=data)
            except requests.exceptions.ConnectionError or requests.exceptions.ConnectTimeout:
                print("[KEYWORD] Connection error.")
                continue

            if response.status_code == 200:
                soup = BeautifulSoup(response.json()['info'])
                products = soup.find_all('div', class_='item-data col-xs-12 col-sm-12')
                foundProduct = False

                for product in products:
                    urlWrapper = product.find('div', class_='img-wrapper')
                    url = urlWrapper.find('a')['href']
                    productName: str = urlWrapper.find('a')['title']

                    if parentThread.firstRun:
                        print(f'[KEYWORDS] Adding firstrun product: {url}')
                        currentProducts.append(url)
                        continue
                    else:
                        if url not in currentProducts and any(keyword in productName.lower().split(' ') for keyword in keywords):
                            currentProducts.append(url)
                            foundProduct = True
                            print(f'[KEYWORDS] Found new product: {url}')

                            while True:
                                pid = GetPIDFromLink(url)

                                print("[KEYWORDS - {}] Getting product info.".format(pid))
                                sizes = []
                                try:
                                    response = requests.post(
                                            'https://www.buzzsneakers.gr/athlitika-papoutsia/',
                                            headers=headers,
                                            data="nbAjax=1&task=getproductdata&productId={}".format(pid)
                                        )
                                except requests.exceptions.ConnectionError:
                                    print("[KEYWORD - {}] Connection error.".format(pid))
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
                                        webhookTitle="New item loaded!",
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

                if parentThread.firstRun:
                    parentThread.firstRun = False

                if not foundProduct:
                    print('[KEYWORDS] No new products found.')
            else:
                print('[KEYWORD] Error fetching data.')
        except Exception as e:
            print('[KEYWORD] Error: {}'.format(e))
            capture_exception(e)
            continue
