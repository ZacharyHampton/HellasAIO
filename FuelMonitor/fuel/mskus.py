from fuel.methods import makeQuery
import time
from db import db
from fuel.keywords import MonitorByKeyword
from fuel.loading import GetMSKUs
from fuel.utils import StockConversion, DontRunKeywords, SizeConversion
from bot.webhook import PingServer, InvalidEmbed
from fuel.size import Size
import traceback
from fuel.product import Product
import os
from dotenv import load_dotenv
import sentry_sdk
from sentry_sdk import capture_exception

load_dotenv()


def MonitorByMSKUs(MSKUs: [str]) -> dict:
    query = """query monitorByMSKU($mskus: [String]!){
  products(sort: {created_at: DESC}, filter: {sku: {in: $mskus}}) {
    items {
      name
      sku
      id
      price {
          regularPrice {
              amount {
                  value
              }
          }
      }
      small_image {
        url
      }
      created_at
      updated_at
      stock_status
      __typename
      ... on ConfigurableProduct {
        variants {
            product {
                id
                created_at
                updated_at
                id 
                name
                sku
                stock_status
                is_raffle_item
            }
        }
      }
    }
  }
}"""

    response = makeQuery(query, {'mskus': MSKUs})
    if response.success:
        return response.data
    else:
        print("Unexpected error: " + response.data)


def MSKUFlowRewrite():
    firstRun = False

    cur = db.cursor()
    previousLength = 0
    sentry_sdk.init(dsn=os.getenv("SENTRY_DSN"))

    while True:
        try:
            time.sleep(1)
            MSKUs = set(GetMSKUs())  #: remove duplicates from msku list
            MSKUs = list(MSKUs)      #: turn back to list
            if len(MSKUs) == 0:      #: if no mskus, skip
                continue

            if len(MSKUs) != previousLength:
                firstRun = True
                previousLength = len(MSKUs)

            data = MonitorByMSKUs(MSKUs)
            if type(data) != dict or data is None:
                print('[msku] Unknown error occurred.')
                continue

            if data.get('errors'):
                print('[msku] Error occurred: ' + data['errors'][0]['message'])
                continue

            if not data['data']['products']['items']:
                print('[msku] No products found.')
                time.sleep(1)
                continue

            for item in data['data']['products']['items']:
                product = Product(
                    mode="msku",
                    productName=item['name'],
                    lastUpdated=item['updated_at'],
                    productId=item['id'],
                    sku=item['sku'],
                    price=float(item['price']['regularPrice']['amount']['value']),
                    imageUrl=item['small_image']['url'],
                    sizes=[],
                )

                row = cur.execute('SELECT sku, lastUpdated FROM msku_data WHERE sku = ?',
                                  (product.sku,)).fetchone()

                #: add every item to db
                if firstRun:
                    if row:
                        product.print('Already in database.')
                    else:
                        product.print('Adding to database')
                        cur.execute('INSERT INTO msku_data (sku, lastUpdated) VALUES (?, ?)', (product.sku, product.lastUpdated))

                    continue

                dbLastUpdated = row[1]
                if dbLastUpdated != product.lastUpdated or row is None:
                    product.print('Product has updated.')
                    product.sizes = [
                        Size(
                            product.sku,
                            SizeConversion(variant['product']['sku']),
                            variant['product']['stock_status'] == "IN_STOCK"
                        ) for variant in item['variants']
                    ]

                    if row is not None:
                        cur.execute('UPDATE msku_data SET lastUpdated = ? WHERE sku = ?', (product.lastUpdated, product.sku))
                    else:
                        cur.execute('INSERT INTO msku_data (sku, lastUpdated) VALUES (?, ?)',
                                    (product.sku, product.lastUpdated))

                    try:
                        PingServer(
                            webhookTitle="Product Update!",
                            productName=product.productName,
                            productPrice=product.price,
                            productMSKU=product.sku,
                            productSizes=product.sizes,
                            productImage=product.imageUrl,
                            productId=product.productId,
                            sizeTitle="Sizes"
                        )
                    except InvalidEmbed as e:
                        capture_exception(e)
                        continue

                else:
                    product.print('Product has not updated.')

            if firstRun:
                firstRun = False
        except Exception as e:
            capture_exception(e)
            continue

