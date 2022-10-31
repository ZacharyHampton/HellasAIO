from discord_webhook import DiscordWebhook, DiscordEmbed
from datetime import datetime
from dotenv import load_dotenv
from buzz.types import Size
import os
import json
import urllib.parse

load_dotenv()


def webhook_time_stmap():
    return datetime.now().strftime("%H:%M:%S")


def PingImageMSKU(MSKU: str, imageURL: str):
    webhook = DiscordWebhook(
        url="discordwebhook",
        rate_limit_retry=True)
    embed = DiscordEmbed(title="Product has been loaded!", description='**{}**'.format(MSKU), color=2524623)
    embed.set_thumbnail(url=imageURL)
    embed.add_embed_field(name="Quicktask Link",
                          value='[Link](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId=3&size=random)'.format(
                              MSKU), inline=True)
    embed.set_footer(text=f"HellasAIO | {webhook_time_stmap()}")
    webhook.add_embed(embed)
    response = webhook.execute()
    if "<Response [405]>" in str(response):
        print("[ERROR] Webhook Incorrect")
    else:
        print("[INFO] Webhook Sent")


def PingServer(webhookTitle: str, productName: str, productMSKU: str, productSizes: list[Size], productImage: str,
               sizeTitle: str, productId: str, productLink: str, includePrice: bool = True, productPrice: float = 0.0, quantity: int = None):
    webhook = DiscordWebhook(url=os.getenv('DISCORD_WEBHOOK'), rate_limit_retry=True)
    webhook2 = DiscordWebhook(url=os.getenv('DISCORD_WEBHOOK2'), rate_limit_retry=True)
    embed = DiscordEmbed(title=webhookTitle, description='**[{}]({})**'.format(productName[:255], productLink),
                         color=2524623)
    embed.set_thumbnail(url=productImage)
    embed.add_embed_field(name="MSKU", value=productMSKU, inline=True)
    embed.add_embed_field(name="PID", value=productId, inline=True)

    if includePrice:
        embed.add_embed_field(name="Price", value="€{:.2f}".format(productPrice), inline=True)

    if productSizes:
        if len('\n'.join([x.webhookText for x in productSizes])) > 1024:
            embed.add_embed_field(name=sizeTitle,
                                  value='\n'.join([x.webhookText for x in productSizes[:len(productSizes) // 2]]),
                                  inline=False)
            embed.add_embed_field(name=sizeTitle,
                                  value='\n'.join([x.webhookText for x in productSizes[len(productSizes) // 2:]]),
                                  inline=True)
        else:
            embed.add_embed_field(name=sizeTitle, value='\n'.join([x.webhookText for x in productSizes]), inline=True)

    if quantity:
        embed.add_embed_field(name="Stock", value=str(quantity), inline=True)

    embed.add_embed_field(name="Quicktask Link",
                          value='[Link](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId=3&size=random)'.format(
                              productMSKU, inline=True))

    embed.set_footer(text=f"HellasAIO | {webhook_time_stmap()}")
    webhook.add_embed(embed)
    response = webhook.execute()
    if "<Response [405]>" in str(response):
        print("[ERROR] Webhook Incorrect")
    else:
        print("[INFO] Webhook Sent")

    embed.set_footer(text=f"Olympus Notify | {webhook_time_stmap()}")
    webhook2.add_embed(embed)
    response = webhook2.execute()
    if "<Response [405]>" in str(response):
        print("[ERROR] Webhook Incorrect")
    else:
        print("[INFO] Webhook Sent")
