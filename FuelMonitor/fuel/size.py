class Size:
    def __init__(self, msku: str, size: str, isInStock: bool):
        self.msku = msku
        self.size: str = size
        self.isInStock: bool = isInStock
        self.webhookText = self.createWebhookText()

    def createWebhookText(self):
        if self.isInStock:
            return '[{}](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId=1&size={}) :green_circle:'.format(self.size, self.msku, self.size)
        else:
            return '[{}](https://quicktask.hellasaio.com/quicktask?product_id={}&siteId=1&size={}) :red_circle:'.format(
                self.size, self.msku, self.size)
