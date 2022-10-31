import random

import requests
import json
from dataclasses import dataclass
import fuel.loading
from fuel.cloudflare import getCloudflareCookies

proxies = fuel.loading.LoadProxies()
CloudflareCookieDict: dict | None = None


@dataclass
class QueryResponse:
    success: bool
    data: dict | str


def makeQuery(query: str, variables: dict) -> QueryResponse:
    global CloudflareCookieDict
    proxy = random.choice(proxies)

    if not CloudflareCookieDict:
        CloudflareCookieDict = getCloudflareCookies(proxy)
        return makeQuery(query, variables)

    """"Sec-Ch-Ua": '"-Not.A/Brand";v="8", "Chromium";v="102"',
            "Sec-Ch-Ua-Mobile": '?0',
            "Sec-Ch-Ua-Platform": '"Windows"',
            "Upgrade-Insecure-Requests": '1',
            "Accept": '*/*',
            "Sec-Fetch-Site": 'same-origin',
            "Sec-Fetch-Mode": 'navigate',
            "Sec-Fetch-Dest": 'document',
            "Accept-Encoding": 'gzip, deflate, br',
            "Connection": 'keep-alive',
            "Accept-Language": "en-US,en;q=0.9",
            "Origin": 'https://www.fuel.com.gr',
            "Cache-Control": 'max-age=0',"""

    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36",
        'Content-Type': 'application/json'
    }

    cookies = CloudflareCookieDict
    data = {'query': query, 'variables': variables}

    try:
        response = requests.post('https://www.fuel.com.gr/el/graphql', json=data, headers=headers,
                                 cookies=cookies
                                 # proxies=proxy
                                 )
    except requests.exceptions.ProxyError:
        print('Proxy failed to make request, retrying...')
        return makeQuery(query, variables)
    except requests.exceptions.ChunkedEncodingError:
        print('Request invalid, retrying...')
        return makeQuery(query, variables)
    except requests.exceptions.SSLError:
        print('SSL error, retrying...')
        return makeQuery(query, variables)

    if (response.status_code == 403 and "1020" in response.text) or (
            response.status_code == 503 and "jschal_js" in response.text):
        # print('[QUERY] Response text,', response.text)
        print('[QUERY] Cloudflare detected.')
        CloudflareCookieDict = getCloudflareCookies(proxy)
        return makeQuery(query, variables)

    if response.status_code == 200:
        return QueryResponse(success=True, data=response.json())
    else:
        return QueryResponse(success=False, data=response.text)
