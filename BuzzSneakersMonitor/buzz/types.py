from dataclasses import dataclass
import threading
import urllib.parse


class Size:
    def __init__(self, pMSKU: str, stock: int, size: str, instock: bool):
        self.msku = pMSKU
        self.stock: int = stock
        self.size: str = size
        self.instock: bool = instock
        self.webhookText = self.createWebhookText()

    def createWebhookText(self):
        if self.instock:
            return '[{}](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId=3&size={}) :green_circle:'.format(urllib.parse.quote(self.size), self.msku, urllib.parse.quote(self.size))
        else:
            return '[{}](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId=3&size={}) :red_circle:'.format(
                urllib.parse.quote(self.size), self.msku, urllib.parse.quote(self.size))


class Thread:
    def __init__(self, flow, pid):
        self.pid = pid
        self.flow = flow
        self.stop = False
        self.thread = threading.Thread(target=self.flow, args=(self.pid, self))
        self.firstRun = True

    def start(self):
        self.thread.start()

    def Stop(self):
        self.stop = True
