from discord_webhook import DiscordWebhook, DiscordEmbed
from datetime import datetime
from dotenv import load_dotenv
from fuel.size import Size
import os
import json

load_dotenv()


class InvalidEmbed(Exception):
    pass


def webhook_time_stmap():
    return datetime.now().strftime("%H:%M:%S")


def PingServer(webhookTitle: str, productName: str, productMSKU: str, productSizes: list[Size], productImage: str, sizeTitle: str, productId: int, includePrice: bool = True, productPrice: float = 0.0):
    webhook = DiscordWebhook(url=os.getenv('DISCORD_WEBHOOK'), rate_limit_retry=True)
    webhook2 = DiscordWebhook(url=os.getenv('DISCORD_WEBHOOK2'), rate_limit_retry=True)
    embed = DiscordEmbed(title=webhookTitle, description='**[{}]({})**'.format(productName[:255], 'https://www.fuel.com.gr/en/catalog/product/view/id/{}'.format(productId)), color=2524623)
    embed.set_thumbnail(url=productImage)
    embed.add_embed_field(name="MSKU", value=productMSKU, inline=True)

    if includePrice:
        embed.add_embed_field(name="Price", value="€{:.2f}".format(productPrice), inline=True)

    if productSizes:
        if len('\n'.join([x.webhookText for x in productSizes])) > 1024:
            embed.add_embed_field(name=sizeTitle, value='\n'.join([x.webhookText for x in productSizes[:len(productSizes) // 2]]), inline=False)
            embed.add_embed_field(name=sizeTitle, value='\n'.join([x.webhookText for x in productSizes[len(productSizes) // 2:]]), inline=True)
        else:
            embed.add_embed_field(name=sizeTitle, value='\n'.join([x.webhookText for x in productSizes]), inline=True)

    embed.add_embed_field(name="Quicktask Link", value='[Link](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId=1&size=random)'.format(productMSKU), inline=True)

    embed.set_footer(text=f"HellasAIO | {webhook_time_stmap()}")
    webhook.add_embed(embed)
    response = webhook.execute()
    if response.status_code == 405:
        print("[ERROR] Webhook Incorrect")
    elif response.status_code == 400:
        print("[ERROR] Webhook Bad Request, invalid embed.")
        raise InvalidEmbed("Invalid embed.", embed)
    else:
        print("[INFO] Webhook Sent")

    embed.set_footer(text=f"Olympus Notify | {webhook_time_stmap()}")
    webhook2.add_embed(embed)
    response = webhook2.execute()
    if response.status_code == 405:
        print("[ERROR] Webhook Incorrect")
    elif response.status_code == 400:
        print("[ERROR] Webhook Bad Request, invalid embed.")
        raise InvalidEmbed("Invalid embed.", embed)
    else:
        print("[INFO] Webhook Sent")
