from fastapi import APIRouter, Depends, HTTPException, status, Response
from dotenv import load_dotenv
import os
from datetime import datetime, timedelta
from pydantic import BaseModel
from fastapi_jwt_auth import AuthJWT
import aiohttp

load_dotenv()
router = APIRouter()
ACCESS_TOKEN_EXPIRE_MINUTES = 10080


class Settings(BaseModel):
    authjwt_secret_key: str = os.getenv("QT_JWT_SECRET")
    authjwt_token_location = {"cookies"}
    authjwt_cookie_domain = "hellasaio.com"
    authjwt_cookie_secure = True


@AuthJWT.load_config
def get_config():
    return Settings()


class QuicktaskAuthRequest(BaseModel):
    code: str
    key: str


class QuicktaskAuthResponse(BaseModel):
    success: bool
    message: str | None = None
    access_token: str | None
    token_type: str | None
    expires_in: int | None


@router.post("/api/quicktask/auth", response_model=QuicktaskAuthResponse)
async def quicktaskauth(data: QuicktaskAuthRequest, response: Response, Authorize: AuthJWT = Depends()):
    async with aiohttp.ClientSession(headers={'Content-Type': 'application/x-www-form-urlencoded'}) as session:
        async with session.post("https://discord.com/api/oauth2/token", data={
            "client_id": os.getenv("DISCORD_CLIENT_ID"),
            "client_secret": os.getenv("DISCORD_CLIENT_SECRET"),
            "grant_type": "authorization_code",
            "code": data.code,
            "redirect_uri": os.getenv("DISCORD_REDIRECT_URL")
        }) as resp:
            if resp.status != 200:
                response.status_code = status.HTTP_401_UNAUTHORIZED
                return QuicktaskAuthResponse(success=False, message="Invalid code")

            jResp = await resp.json()

    async with aiohttp.ClientSession(headers={'Authorization': "Bearer {}".format(jResp['access_token'])}) as session:
        async with session.get("https://discord.com/api/users/@me") as resp:
            if resp.status != 200:
                jResp.status_code = status.HTTP_401_UNAUTHORIZED
                return QuicktaskAuthResponse(success=False, message="Failed to get discord account")

            jResp = await resp.json()
            discordIdFromDiscord = jResp['id']

    async with aiohttp.ClientSession(
            headers={
                'Accept': 'application/json',
                'Authorization': 'Bearer {}'.format((os.getenv('WHOP_BEARER')))
            }) as session:
        async with session.get('https://api.whop.com/api/v1/licenses/{}'.format(data.key)) as resp:
            jResp = await resp.json()
            if jResp.get('message'):
                if jResp['message'] == "Please reset your key to use on a new machine":
                    response.status_code = status.HTTP_401_UNAUTHORIZED
                    return QuicktaskAuthResponse(success=False, message="HWID does not match current computer's HWID.")

                if jResp['message'] == "Not found":
                    response.status_code = status.HTTP_401_UNAUTHORIZED
                    return QuicktaskAuthResponse(success=False, message="License key not found.")

                if jResp['message'] == "Please confirm your API token":
                    response.status_code = status.HTTP_401_UNAUTHORIZED
                    return QuicktaskAuthResponse(success=False, message="Backend API Key Error.")

                response.status_code = status.HTTP_503_SERVICE_UNAVAILABLE
                return QuicktaskAuthResponse(success=False, message="Unknown error.")

            if jResp.get('banned'):
                response.status_code = status.HTTP_401_UNAUTHORIZED
                return QuicktaskAuthResponse(success=False, message="License key is banned.")

            if jResp.get('is_scammer'):
                response.status_code = status.HTTP_401_UNAUTHORIZED
                return QuicktaskAuthResponse(success=False, message="License key is marked as a scammer.")

            if not jResp.get('valid'):
                response.status_code = status.HTTP_401_UNAUTHORIZED
                return QuicktaskAuthResponse(success=False, message="License key is invalid.")

            if any(
                    [
                        jResp['key_status'] == 'approved',
                        jResp['key_status'] == 'listed'
                    ]
            ) and any(
                [
                    jResp['subscription_status'] == 'completed',
                    jResp['subscription_status'] == 'active',
                    jResp['subscription_status'] == 'trialing'
                ]
            ):
                if str(jResp['discord']['discord_account_id']) != discordIdFromDiscord:
                    response.status_code = status.HTTP_401_UNAUTHORIZED
                    return QuicktaskAuthResponse(success=False, message="Discord ID does not match key's Discord ID.")

                if jResp['metadata'].get('HWID') is None:
                    response.status_code = status.HTTP_401_UNAUTHORIZED
                    return QuicktaskAuthResponse(success=False, message="No HWID binded to key.")

                access_token = Authorize.create_access_token(
                    subject=data.key,
                    expires_time=ACCESS_TOKEN_EXPIRE_MINUTES * 60,
                    user_claims={'discord_id': discordIdFromDiscord}
                )

                Authorize.set_access_cookies(access_token, response, max_age=ACCESS_TOKEN_EXPIRE_MINUTES * 60)
                response.set_cookie("accessBool", "true", max_age=ACCESS_TOKEN_EXPIRE_MINUTES * 60, domain="hellasaio.com")

                return QuicktaskAuthResponse(
                    success=True,
                    access_token=access_token,
                    token_type="Bearer",
                    expires_in=ACCESS_TOKEN_EXPIRE_MINUTES
                )
            else:
                response.status_code = status.HTTP_401_UNAUTHORIZED
                return QuicktaskAuthResponse(success=False, message="License key is not approved.")
