import os
import requests
from fastapi import APIRouter
from dotenv import load_dotenv
from datetime import datetime, timedelta
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from jose import JWTError, jwt
from pydantic import BaseModel

load_dotenv()
router = APIRouter()
SECRET_KEY = os.getenv("JWT_SECRET")
ALGORITHM = "HS256"
ACCESS_TOKEN_EXPIRE_MINUTES = 1
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="auth")


class TokenData(BaseModel):
    key: str
    discordId: int


class Login(BaseModel):
    licenseKey: str
    HWID: str  #: sha256(Disk Serials (sep by comma) + Computer Name + Running User)


class LoginResponse(BaseModel):
    success: bool
    message: str | None = None
    access_token: str | None = None
    token_type: str | None = None


def create_access_token(data: dict):
    to_encode = data.copy()
    expire = datetime.utcnow() + timedelta(minutes=ACCESS_TOKEN_EXPIRE_MINUTES)
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt


async def get_current_user(token: str = Depends(oauth2_scheme), sKey: str = SECRET_KEY):
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )
    try:
        payload = jwt.decode(token, sKey, algorithms=[ALGORITHM])
        key: str = payload.get("sub")
        if key is None:
            raise credentials_exception
        token_data = TokenData(key=key, discordId=payload.get("discord_id"))
    except JWTError:
        raise credentials_exception
    return token_data


@router.post("/api/auth", response_model=LoginResponse)
async def auth(login: Login):
    response = requests.post(
        'https://api.whop.com/api/v1/licenses/{}/validate'.format(login.licenseKey),
        json={'metadata': {'HWID': login.HWID}},
        headers={
            'Accept': 'application/json',
            'Authorization': 'Bearer {}'.format((os.getenv('WHOP_BEARER')))
        })

    try:
        response = response.json()
    except requests.exceptions.JSONDecodeError:
        return LoginResponse(success=False, message="Backend login error.")

    if response.get('message'):
        if response['message'] == "Please reset your key to use on a new machine":
            return LoginResponse(success=False, message="HWID does not match current computer's HWID.")

        if response['message'] == "Not found":
            return LoginResponse(success=False, message="License key not found.")

        if response['message'] == "Please confirm your API token":
            return LoginResponse(success=False, message="Backend API Key Error")

        return LoginResponse(success=False, message="Unknown error.")

    if response.get('banned'):
        return LoginResponse(success=False, message="License key is banned.")

    if response.get('is_scammer'):
        return LoginResponse(success=False, message="License key is marked as a scammer.")

    if not response.get('valid'):
        return LoginResponse(success=False, message="License key is invalid.")

    if any(
            [
                response['key_status'] == 'approved',
                response['key_status'] == 'listed'
            ]
    ) and any(
        [
            response['subscription_status'] == 'completed',
            response['subscription_status'] == 'active',
            response['subscription_status'] == 'trialing'
        ]
    ):
        return LoginResponse(success=True, token_type="bearer", access_token=create_access_token(
            data={"sub": response['key'], "discord_id": response['discord']['discord_account_id']},
        ))
    else:
        return LoginResponse(success=False, message="License key is not approved.")
