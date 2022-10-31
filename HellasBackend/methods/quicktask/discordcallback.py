from fastapi import APIRouter
from fastapi.responses import RedirectResponse

router = APIRouter()


@router.get("/api/discord/callback", response_class=RedirectResponse, status_code=302)
async def discord_callback(code: str):
    return RedirectResponse("https://quicktask.hellasaio.com/key?code=" + code)
