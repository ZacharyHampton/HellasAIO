import requests
from fastapi import APIRouter, Depends, Header
from dotenv import load_dotenv
from methods.auth import TokenData, get_current_user
from pydantic import BaseModel
import helheim
import cloudscraper
from urllib.parse import urlparse


load_dotenv()
router = APIRouter()

SiteIDToSiteURL = {
    0: "https://www.theathletesfoot.gr",
    1: "https://www.fuel.com.gr",
    2: "https://www.slamdunk.gr/",
    3: "https://www.buzzsneakers.gr/",
    4: "https://europesports.eu/"
}


class CookieResponse(BaseModel):
    cookies: dict | None = None
    success: bool
    message: str | None = None


def injection(session, response):
    if helheim.isChallenge(session, response):
        return helheim.solve(session, response)
    else:
        return response


@router.get('/api/cloudflare/{siteId}', response_model=CookieResponse)
async def getCloudflareCookies(
        siteId: int,
        current_user: TokenData = Depends(get_current_user),
        user_agent: str | None = Header(default=None),
        x_proxy: str | None = Header(default=None)
):
    if x_proxy is None:
        return CookieResponse(success=False, message="No proxy provided.")

    if siteId:
        if siteId != 1:
            return CookieResponse(success=False, message="This site is not supported yet.")

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

    if user_agent:
        session.headers['User-Agent'] = user_agent

    if x_proxy:
        """pUrl = urlparse(x_proxy)

        session.proxies = {
            'http': "http://{}:{}@{}:{}".format(pUrl.username, pUrl.password, pUrl.hostname, pUrl.port),
            'https': "https://{}:{}@{}:{}".format(pUrl.username, pUrl.password, pUrl.hostname, pUrl.port)
        }"""

        session.proxies = {
            'http': x_proxy,
            'https': x_proxy
        }

    try:
        response = session.get(SiteIDToSiteURL[siteId])
    except helheim.exceptions.HelheimRuntimeError or requests.exceptions.ProxyError or requests.exceptions.SSLError or requests.exceptions.ConnectionError:
        return CookieResponse(success=False, message="Proxy error. Cannot connect to proxy.")
    except requests.exceptions.ChunkedEncodingError:
        return CookieResponse(success=False, message="Chunked encoding error.")

    if response.status_code == 200:
        return CookieResponse(success=True, cookies=session.cookies.get_dict())
    else:
        return CookieResponse(success=False, message="Something went wrong. ({})".format(response.status_code))
