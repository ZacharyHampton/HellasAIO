from fastapi import APIRouter, Cookie, Response, status
from methods.quicktask import qtWebsocket
from methods.auth import get_current_user
import os
from dotenv import load_dotenv
from pydantic import BaseModel

load_dotenv()
router = APIRouter()
SECRET_KEY = os.getenv("QT_JWT_SECRET")


class QTResponse(BaseModel):
    success: bool
    message: str


#: will be vuln to csrf
@router.get('/api/quicktask/start')
async def startQuicktask(siteId: str, product_id: str, size: str, response: Response,
                         access_token_cookie: str = Cookie()):
    user = await get_current_user(access_token_cookie, SECRET_KEY)

    if qtWebsocket.manager.get_connections(user.key) is None:
        response.status_code = status.HTTP_400_BAD_REQUEST
        return QTResponse(success=False, message="No active client running.")

    await qtWebsocket.manager.send_json(user.key, {"siteId": siteId, "product_id": product_id, "size": size})
    return QTResponse(success=True, message="Sent to client.")
