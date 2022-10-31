import time
import requests
import threading
from db import db
from buzz import utils, webhook
from buzz.types import Size, Thread
import sentry_sdk
import os
from dotenv import load_dotenv
from sentry_sdk import capture_exception

load_dotenv()
sentry_sdk.init(dsn=os.getenv('SENTRY_DSN'))


def BuzzPIDFlow(PID: str, parentThread: Thread):
    print('[{}] Starting thread.'.format(PID))

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

    cur = db.cursor()

    while not parentThread.stop:
        try:
            time.sleep(1)

            try:
                response = requests.post(
                    'https://www.buzzsneakers.gr/athlitika-papoutsia/',
                    headers=headers,
                    data="nbAjax=1&task=getproductdata&productId={}".format(PID)
                )
            except requests.exceptions.ConnectionError:
                print('[{}] Connection error.'.format(PID))
                continue

            row = cur.execute("SELECT (instock) FROM product_data WHERE id = ?", (PID,)).fetchone()

            if response.status_code == 200:
                rJSON = response.json()
                print("[{}] Successfully fetched data.".format(PID))

                if parentThread.firstRun:
                    if row:
                        parentThread.firstRun = False
                        continue
                    else:
                        cur.execute("INSERT INTO product_data (instock, id) VALUES (?, ?)", (utils.GetStockBool(rJSON), PID))
                        db.commit()
                else:
                    if bool(row[0]) != bool(utils.GetStockBool(rJSON)):
                        cur.execute("UPDATE product_data SET instock = ? WHERE id = ?", (utils.GetStockBool(rJSON), PID))
                        db.commit()

                    if not bool(row[0]) and bool(utils.GetStockBool(rJSON)):
                        print("[{}] Product is now in stock!".format(PID))
                        sizes: list[Size] = []

                        for size in rJSON['sizes']:
                            sizes.append(Size(
                                pMSKU=rJSON['product']['productCode'],
                                stock=int(float(size['quantity'])),
                                size=size['sizeName'],
                                instock=bool(utils.GetStockBoolSize(size))
                            ))

                        webhook.PingServer(
                            webhookTitle="Product Instock!",
                            productName=rJSON['product']['name'],
                            productMSKU=rJSON['product']['productCode'],
                            productSizes=sizes,
                            productImage="https://www.buzzsneakers.gr/{}".format(rJSON['product']['image']),
                            sizeTitle="Sizes",
                            productId=PID,
                            productLink=rJSON['product']['permalinkProduct'],
                            includePrice=False,
                            quantity=int(float(rJSON['product']['quantity']))
                        )
                    else:
                        print('[{}] Product has not changed.'.format(PID))

            else:
                print("[{}] Error fetching data.".format(PID))
        except Exception as e:
            capture_exception(e)
            continue




