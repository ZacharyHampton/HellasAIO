from db import db


def GetPIDs() -> list[str]:
    cur = db.cursor()
    PIDs = []

    for row in cur.execute('SELECT id FROM products'):
        PIDs.append(row[0])

    return PIDs


def GetStockBool(productData: dict) -> int:
    return int(float(productData['product']['quantity']) != 0)


def GetStockBoolSize(sizeData: dict) -> int:
    return int(float(sizeData['quantity']) != 0)


def GetPIDFromLink(link: str) -> str:
    return link.split('/')[-1].split('-')[0]
