from db import db


def GetKeywords() -> list[str]:
    cur = db.cursor()
    Keywords = []

    for row in cur.execute('SELECT keyword FROM keywords'):
        Keywords.append(row[0])

    return Keywords


def GetMSKUs() -> list[str]:
    cur = db.cursor()
    MSKUs = []

    for row in cur.execute('SELECT msku FROM mskus'):
        MSKUs.append(row[0])

    return MSKUs


def ConvertProxy(proxy: str) -> dict:
    proxyParts = proxy.strip().split(':')
    proxy = '{}:{}@{}:{}'.format(proxyParts[2], proxyParts[3], proxyParts[0], proxyParts[1])
    return {'http': 'http://' + proxy, 'https': 'http://' + proxy}


def LoadProxies() -> list[dict]:
    with open('proxies.txt', 'r') as f:
        proxies = f.readlines()

    return [ConvertProxy(proxy) for proxy in proxies]
