from fastapi import APIRouter, Depends, Request, status, HTTPException, BackgroundTasks
from methods.deployment.compiling import BuildNewRelease
from dotenv import load_dotenv
import os
import hmac
from pydantic import BaseModel

load_dotenv()
router = APIRouter()
WEBHOOK_SECRET = os.getenv("GITHUB_WEBHOOK_SECRET")


class Response(BaseModel):
    success: str
    message: str


async def validateRequest(
        request: Request
):
    requestBody = await request.body()
    hmacObject = hmac.new(WEBHOOK_SECRET.encode(), requestBody)

    valid = hmac.compare_digest(
        hmacObject.hexdigest(),
        request.headers.get("X-Hub-Signature-256")
    )

    if not valid:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED)


@router.post('/api/github/release', dependencies=[Depends(validateRequest)], response_model=Response)
async def github_release_webhook(
        request: Request,
        background_tasks: BackgroundTasks
):
    responseJSON = await request.json()

    if responseJSON['action'] != "released":
        return Response(success=False, message="Release is not action 'released'.")

    description = responseJSON['release']['body']
    version = responseJSON['release']['name']

    background_tasks.add_task(BuildNewRelease, description, version)
    return Response(success=True, message="Starting build of new release.")


