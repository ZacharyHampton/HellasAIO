import helheim
import cloudscraper
from urllib.parse import urlparse
from dotenv import load_dotenv
import os

load_dotenv()
helheim.auth(os.getenv('HELHEIM_KEY'))


def injection(session, response):
    if helheim.isChallenge(session, response):
        return helheim.solve(session, response)
    else:
        return response


def getCloudflareCookies(proxy: dict) -> dict:
    session = cloudscraper.create_scraper(
        browser={
            'browser': 'chrome',  # we want a chrome user-agent
            'mobile': False,  # pretend to be a desktop by disabling mobile user-agents
            'platform': 'windows'  # pretend to be 'windows' or 'darwin' by only giving this type of OS for user-agents
        },
        requestPostHook=injection,
        captcha={
            'provider': 'vanaheim',
        }
    )
    helheim.wokou(session)

    session.headers[
        'User-Agent'] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " \
                        "Chrome/103.0.0.0 Safari/537.36"

    #session.proxies = proxy

    try:
        response = session.get("https://www.fuel.com.gr")
    except helheim.exceptions.HelheimRuntimeError as e:
        print('[CLOUDFLARE]', e)
        print("[CLOUDFLARE] Proxy error. Cannot connect to proxy.")
        return getCloudflareCookies(proxy)

    if response.status_code == 200:
        print("[CLOUDFLARE] Successfully generated cookies.")
        return session.cookies.get_dict()
    else:
        print("[CLOUDFLARE] Something went wrong. ({})".format(response.status_code))
