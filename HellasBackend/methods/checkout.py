import os
from fastapi import APIRouter, Depends, BackgroundTasks
from methods.auth import TokenData, get_current_user
from pydantic import BaseModel
import requests
import time
from dotenv import load_dotenv
import urllib.parse

load_dotenv()
router = APIRouter()

SiteConversion = {
    "athletesfoot": 0,
    "fuel": 1,
    "slamdunk": 2,
    "buzzsneakers": 3,
    "europesports": 4,
}


class CheckoutResponse(BaseModel):
    success: bool
    message: str | None = None


class Checkout(BaseModel):
    price: float | None = None,
    product_name: str | None = None,
    product_msku: str | None = None
    mode: str | None = None,
    checkout_ms: int | None = None
    size: str | None = None,
    status: str | None = None,
    website: str | None = None,
    image_url: str | None = None,
    allow_public: bool | None = None


def _remove_none_values(d: dict):
    return {k: v for k, v in d.items() if v is not None}


def send_webhook(webhookPost: dict):
    while True:
        response = requests.post(os.getenv('GROUP_WEBHOOK'), json=webhookPost,
                                 headers={'content-type': 'application/json'})
        if response.status_code == 429:
            time.sleep(response.json()['retry_after'])
        elif response.status_code == 204:
            break
        else:
            print(response.text)
            print(response.status_code)
            break


def createSuccessBody(checkout: Checkout, user: TokenData):
    return {
        "content": None,
        "embeds": [
            {
                "title": "**Successful Checkout!**",
                "description": checkout.product_name,
                "color": 2524623,
                "fields": [
                    {
                        "name": "MSKU",
                        "value": checkout.product_msku,
                        "inline": True
                    },
                    {
                        "name": "Mode",
                        "value": checkout.mode,
                        "inline": True
                    },
                    {
                        "name": "Size",
                        "value": "[{}](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId={}&size={})".format(checkout.size, checkout.product_msku, SiteConversion[checkout.website.lower()], urllib.parse.quote(checkout.size)),
                        "inline": True
                    },
                    {
                        "name": "Checkout Time",
                        "value": "{}ms".format(checkout.checkout_ms),
                        "inline": True
                    },
                    {
                        "name": "Price",
                        "value": "€{:.2f}".format(checkout.price),
                        "inline": True
                    },
                    {
                        "name": "Store",
                        "value": checkout.website,
                        "inline": True
                    },
                    {
                        "name": "Quicktask Link",
                        "value": "[Link](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId={}&size=random)".format(checkout.product_msku, SiteConversion[checkout.website.lower()]),
                        "inline": True
                    },
                    {
                        "name": "User",
                        "value": "<@{}>".format(user.discordId),
                        "inline": True
                    }
                ],
                "thumbnail": {
                    "url": checkout.image_url
                }
            }
        ],
        "attachments": []
    }


@router.post("/api/checkout")
async def checkoutLog(checkout: Checkout, background_tasks: BackgroundTasks,
                      current_user: TokenData = Depends(get_current_user)):
    requestBody = {
        "key": current_user.key,
        "price": checkout.price,
        "product_name": checkout.product_name,
        "size": checkout.size,
        "status": checkout.status,
        "website": checkout.website,
        "image_url": checkout.image_url,
    }
    requestBody = _remove_none_values(requestBody)
    response = requests.post(
        'https://api.whop.com/api/v1/checkout_logs',
        headers={
            'Accept': 'application/json',
            'Authorization': 'Bearer {}'.format((os.getenv('WHOP_BEARER')))
        },
        json=requestBody)

    if checkout.status == "success" and checkout.allow_public:
        background_tasks.add_task(send_webhook, createSuccessBody(checkout, current_user))

    return CheckoutResponse(success=response.status_code == 200)
