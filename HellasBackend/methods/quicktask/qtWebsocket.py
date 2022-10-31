import json
from fastapi import WebSocket, Depends, HTTPException, Query, status, WebSocketDisconnect
from fastapi import APIRouter
from dotenv import load_dotenv
from methods.auth import TokenData, get_current_user
from pydantic import BaseModel


class ConnectionManager:
    def __init__(self):
        self.active_connections: dict[str][list[WebSocket]] = {}

    async def connect(self, ws: WebSocket, token: str) -> TokenData | None:
        await ws.accept()

        try:
            user = await get_current_user(token)
        except HTTPException as e:
            if e.status_code == 401:
                await ws.close(code=status.WS_1008_POLICY_VIOLATION)
                return None
            else:
                await ws.send_json(LoginResponse(success=False, error=e.status_code).dict())
                await ws.close()
                return None

        await manager.add(user.key, ws)
        return user

    async def add(self, key: str, ws: WebSocket):
        if self.active_connections.get(key) is None:
            self.active_connections[key] = [ws]
        else:
            self.active_connections[key].append(ws)

    def get_connections(self, key: str):
        return self.active_connections.get(key)

    async def disconnect(self, key: str, ws: WebSocket):
        connection = self.active_connections.get(key)
        if connection is None:
            return

        self.active_connections[key].remove(ws)
        if not self.active_connections[key]:
            del self.active_connections[key]

    async def send_json(self, key: str, data: dict):
        for connection in self.get_connections(key):
            await connection.send_json(data)

    """async def broadcast(self, data: dict):
        for connection in self.active_connections.values():
            for ws in connection:
                await ws.send_json(data)"""


class LoginResponse(BaseModel):
    success: bool
    error: int | None = None
    key: str | None = None


async def get_token(
        ws: WebSocket,
        token: str | None = Query(default=None),
):
    if token is None:
        await ws.close(code=status.WS_1008_POLICY_VIOLATION)
    return token

load_dotenv()
router = APIRouter()
manager = ConnectionManager()


@router.websocket('/api/ws')
async def websocket(ws: WebSocket, token: str = Depends(get_token)):
    user = await manager.connect(ws, token)
    if user is None:
        return

    await ws.send_json(LoginResponse(success=True, key=user.key).dict())

    try:
        while user.key in manager.active_connections:
            try:
                data = await ws.receive_text()
            except WebSocketDisconnect:
                await manager.disconnect(key=user.key, ws=ws)
                break
    except WebSocketDisconnect:
        await manager.disconnect(key=user.key, ws=ws)
    except json.JSONDecodeError:
        await manager.send_json(user.key, {"error": "Invalid JSON"})

