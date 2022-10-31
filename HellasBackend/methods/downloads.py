import os
from fastapi import APIRouter, Depends, Response, status, Header
from methods.auth import TokenData, get_current_user
from pydantic import BaseModel
import clients
import requests
import time
from dotenv import load_dotenv

load_dotenv()
router = APIRouter()
db = clients.MongoClient.HellasAIO


class DownloadData(BaseModel):
    url: str
    checksum: str


class Downloads(BaseModel):
    windows: DownloadData
    macos_intel: DownloadData
    macos_m1: DownloadData


class DownloadsResponse(BaseModel):
    success: bool
    version: str | None = None
    downloads: Downloads | None = None


class SetDownloadsRequest(BaseModel):
    version: str
    downloads: Downloads


@router.get("/api/latest", response_model=DownloadsResponse)
def getLatestDownloads():
    latestDownloadData = db.Downloads.find().sort("_id", -1)
    return DownloadsResponse(
        success=True,
        version=latestDownloadData[0]['version'],
        downloads=Downloads(
            windows=DownloadData(
                url=latestDownloadData[0]['downloads']['windows']['url'],
                checksum=latestDownloadData[0]['downloads']['windows']['checksum']
            ),
            macos_intel=DownloadData(
                url=latestDownloadData[0]['downloads']['macos_intel']['url'],
                checksum=latestDownloadData[0]['downloads']['macos_intel']['checksum']
            ),
            macos_m1=DownloadData(
                url=latestDownloadData[0]['downloads']['macos_m1']['url'],
                checksum=latestDownloadData[0]['downloads']['macos_m1']['checksum']
            )
        )
    )


@router.post("/api/downloads", response_model=DownloadsResponse, status_code=status.HTTP_200_OK)
def setDownloads(
        data: SetDownloadsRequest,
        response: Response,
        current_user: TokenData = Depends(get_current_user),
        x_password: str = Header()
):
    if current_user.discordId != int(os.getenv("OWNER_ID")) and x_password != os.getenv("UPLOAD_PASSWORD"):
        response.status_code = status.HTTP_401_UNAUTHORIZED
        return DownloadsResponse(success=False)

    db.Downloads.insert_one(data.dict())
    return DownloadsResponse(success=True, version=data.version, downloads=data.downloads)
