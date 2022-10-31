from dataclasses import dataclass
from fuel import size


@dataclass
class Product:
    mode: str
    productName: str
    productId: int
    sku: str
    lastUpdated: str
    price: float
    imageUrl: str
    sizes: list[size.Size]

    def print(self, *message):
        print('[{}] [{}] {}'.format(self.mode, self.sku, *message))



