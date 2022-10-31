import time
from db import db
from fuel.mskus import MonitorByMSKUs
from fuel.keywords import MonitorByKeyword
from fuel.loading import GetMSKUs
from fuel.utils import StockConversion, DontRunKeywords
from bot.webhook import PingServer, InvalidEmbed
import fuel.size
import traceback
import os
from dotenv import load_dotenv
import sentry_sdk
from sentry_sdk import capture_exception

load_dotenv()


def MSKUFlow():
    firstRun = True
    sentry_sdk.init(dsn=os.getenv("SENTRY_DSN"))

    try:
        cur = db.cursor()

        while True:
            time.sleep(1)
            MSKUs = set(GetMSKUs())
            MSKUs = list(MSKUs)
            if len(MSKUs) == 0:
                continue

            data = MonitorByMSKUs(MSKUs)
            if not data['data']['products']['items']:
                print('[msku mode] No products found.')
                time.sleep(1)
                continue

            for item in data['data']['products']['items']:
                productName = item['name']
                productId = item['id']
                sku = item['sku']
                price = item['price']['regularPrice']['amount']['value']
                imageUrl = item['small_image']['url']
                sizes: list[fuel.size.Size] = []
                NeedToAlertSizes = []

                print('[{}] [{}] Checking item...'.format(sku, productName))

                if item['stock_status'] == 'IN_STOCK':
                    print('[{}] [{}] Product has variants.'.format(sku, productName))
                    for variant in item['variants']:
                        print('[{}] [{}] Checking variant: {}'.format(sku, productName, variant['product']['sku']))
                        if firstRun:
                            row = cur.execute('SELECT variantsku, instock FROM msku_data WHERE variantsku = ?',
                                              (variant['product']['sku'],)).fetchone()
                            if row is None:
                                print('[{}] [{}] Adding variant to database (first run)...'.format(sku, productName))
                                cur.execute('INSERT INTO msku_data (variantsku, instock) VALUES (?, ?)',
                                            (
                                                variant['product']['sku'],
                                                StockConversion(variant['product']['stock_status'])))
                            else:
                                print('[{}] [{}] Variant already in database (first run)...'.format(sku, productName))

                            continue

                        variantSizeSplit = variant['product']['sku'].split(' ')
                        if len(variantSizeSplit) > 1:
                            sizes.append(fuel.size.Size(sku, variant['product']['sku'].split(' ')[1],
                                                        variant['product']['stock_status'] == "IN_STOCK"))
                        else:
                            sizes.append(fuel.size.Size(sku, variant['product']['sku'].split('-')[-1],
                                                        variant['product']['stock_status'] == "IN_STOCK"))
                        if variant['product']['stock_status'] == 'IN_STOCK':  #: if variant instock
                            row = cur.execute('SELECT variantsku, instock FROM msku_data WHERE variantsku = ?',
                                              (variant['product']['sku'],)).fetchone()
                            if row is not None:  #: if variant in db
                                if row[1] == 1:  #: if in stock
                                    print('[{}] [{}] Variant {} was previously in stock.'.format(sku, productName,
                                                                                                 variant['product'][
                                                                                                     'sku']))
                                    continue

                            #: problem: if variant not in db and is in stock, it will ping the channel; this is okay?
                            print('[{}] [{}] Variant "{}" in stock, pinging...'.format(sku, productName,
                                                                                       variant['product']['name']))
                            if len(variantSizeSplit) > 1:
                                NeedToAlertSizes.append(variant['product']['sku'].split(' ')[1])
                            else:
                                NeedToAlertSizes.append(variant['product']['sku'].split('-')[1])

                            cur.execute('INSERT INTO msku_data (variantsku, instock) VALUES (?, ?)',
                                        (
                                            variant['product']['sku'],
                                            StockConversion(variant['product']['stock_status'])))
                        else:
                            #: oos: if variant out of stock
                            row = cur.execute('SELECT variantsku, instock FROM msku_data WHERE variantsku = ?',
                                              (variant['product']['sku'],)).fetchone()
                            if row is not None:  #: if variant in db
                                if row[1] != 0:  #: if not oos in db
                                    cur.execute('UPDATE msku_data SET instock = ? WHERE variantsku = ?',
                                                (0, variant['product']['sku']))

                else:
                    print('[{}] [{}] Product out of stock.'.format(sku, productName))
                    continue  #: skip if oos

                if len(NeedToAlertSizes) > 0:
                    try:
                        PingServer(
                            webhookTitle="Product In-Stock!",
                            productName=productName,
                            productPrice=float(price),
                            productMSKU=sku,
                            productSizes=sizes,
                            productImage=imageUrl,
                            productId=productId,
                            sizeTitle="Sizes"
                        )
                    except InvalidEmbed as e:
                        capture_exception(e)
                        continue

            if firstRun:
                firstRun = False

            db.commit()
    except Exception as e:
        print(traceback.format_exc())
        capture_exception(e)
        MSKUFlow()


def KeywordFlow(keyword: str):
    cur = db.cursor()
    firstRun = True

    try:
        while True:
            if keyword in DontRunKeywords:
                break

            time.sleep(1)
            data = MonitorByKeyword(keyword)
            print('[{}] Checking keyword...'.format(keyword))

            if data is None or type(data) != dict:
                print('Unknown error occurred.')
                print(data)
                print(type(data))
                continue

            if data.get('errors'):
                print('[{}] Error occurred: '.format(keyword) + data['errors'][0]['message'])
                continue

            if not data['data']['products']['items']:
                print('[{}] No products found.'.format(keyword))
                time.sleep(1)
                continue

            for item in data['data']['products']['items']:
                productName = item['name']
                productId = item['id']
                lastUpdated = item['updated_at']
                sku = item['sku']
                imageUrl = item['small_image']['url']
                sizes = []

                print('[{}] [{}] Checking item...'.format(keyword, productName))
                row = cur.execute('SELECT msku, last_updated FROM keyword_data WHERE msku = ?', (sku,)).fetchone()

                if firstRun:
                    if row is None:
                        print('[{}] [{}] Product not in database, adding (first run)...'.format(keyword, productName))
                        cur.execute('INSERT INTO keyword_data (msku, keyword, last_updated) VALUES (?, ?, ?)',
                                    (sku, keyword, lastUpdated))
                    else:
                        print('[{}] [{}] Product already in database (first run)...'.format(keyword, productName))

                    continue

                #: main flow (not first run)
                if row is None or row[1] != lastUpdated:  #: new item/updated
                    print(row)  # debug
                    print('[{}] [{}] New item, pinging...'.format(keyword, productName))
                    if row is None:
                        cur.execute('INSERT INTO keyword_data (msku, keyword, last_updated) VALUES (?, ?, ?)',
                                    (sku, keyword, lastUpdated))
                    else:
                        cur.execute('UPDATE keyword_data SET last_updated = ? WHERE msku = ?', (lastUpdated, sku))

                    for variant in item.get('variants', []):
                        variantSizeSplit = variant['product']['sku'].split(' ')
                        if len(variantSizeSplit) > 1:
                            sizes.append(fuel.size.Size(sku, variant['product']['sku'].split(' ')[1],
                                                        variant['product']['stock_status'] == "IN_STOCK"))
                        else:
                            sizes.append(fuel.size.Size(sku, variant['product']['sku'].split('-')[-1],
                                                        variant['product']['stock_status'] == "IN_STOCK"))
                    try:
                        PingServer(
                            webhookTitle="Product Update!",
                            productName=productName,
                            productMSKU=sku,
                            productSizes=sizes,
                            productImage=imageUrl,
                            productId=productId,
                            sizeTitle="Sizes",
                            includePrice=False
                        )
                    except InvalidEmbed as e:
                        capture_exception(e)
                        continue


                else:
                    print('[{}] [{}] Item already in database.'.format(keyword, productName))

            firstRun = False

            """rows = cur.execute('SELECT msku FROM keyword_data WHERE keyword = ?', (keyword,)).fetchall()
            wFetchedMSKUs = [x['sku'] for x in data['data']['products']['items']]
            dFetchedMSKUs = [row[0] for row in rows]
            for msku in dFetchedMSKUs:
                if msku not in wFetchedMSKUs:
                    cur.execute('DELETE FROM keyword_data WHERE msku = ?', (msku,))
                    print('[{}] [{}] Item removed from database.'.format(keyword, msku))"""

            db.commit()
    except Exception as e:
        print(traceback.format_exc())
        capture_exception(e)
        KeywordFlow(keyword)
