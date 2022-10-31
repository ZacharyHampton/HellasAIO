import os
from fastapi import FastAPI, Request, Response
from fastapi.responses import JSONResponse
from methods import auth, checkout, downloads, cloudflare
from methods.quicktask import discordauth, discordcallback, qtWebsocket, quicktaskapi
from methods.deployment import webhook
import sentry_sdk
from dotenv import load_dotenv
from fastapi_jwt_auth.exceptions import AuthJWTException
from starlette.middleware.cors import CORSMiddleware
from starlette.middleware import Middleware

origins = [
    "http://localhost:8000",
    "https://api.hellasaio.com"
    "https://quicktask.hellasaio.com"
    "https://*.hellasaio.com"
    "http://*.hellasaio.com"
]

load_dotenv()
App = FastAPI(redoc_url=None, docs_url=None)
App.include_router(auth.router)
App.include_router(checkout.router)
App.include_router(downloads.router)
App.include_router(cloudflare.router)
App.include_router(discordcallback.router)
App.include_router(discordauth.router)
App.include_router(qtWebsocket.router)
App.include_router(quicktaskapi.router)
App.include_router(webhook.router)

# Salt to your taste
ALLOWED_ORIGINS = 'https://quicktask.hellasaio.com'  # or 'foo.com', etc.


# handle CORS preflight requests
@App.options('/{rest_of_path:path}')
async def preflight_handler(request: Request, rest_of_path: str) -> Response:
    response = Response()
    response.headers['Access-Control-Allow-Origin'] = ALLOWED_ORIGINS
    response.headers['Access-Control-Allow-Methods'] = 'POST, GET, DELETE, OPTIONS'
    response.headers['Access-Control-Allow-Headers'] = 'Authorization, Content-Type'
    response.headers['Access-Control-Allow-Credentials'] = 'true'
    return response


# set CORS headers
@App.middleware("http")
async def add_CORS_header(request: Request, call_next):
    response = await call_next(request)
    response.headers['Access-Control-Allow-Origin'] = ALLOWED_ORIGINS
    response.headers['Access-Control-Allow-Methods'] = 'POST, GET, DELETE, OPTIONS'
    response.headers['Access-Control-Allow-Headers'] = 'Authorization, Content-Type'
    response.headers['Access-Control-Allow-Credentials'] = 'true'
    return response


@App.exception_handler(AuthJWTException)
def authjwt_exception_handler(request: Request, exc: AuthJWTException):
    return JSONResponse(
        status_code=exc.status_code,
        content={"detail": exc.message}
    )


sentry_sdk.init(
    dsn=os.getenv('SENTRY_DSN'),
    traces_sample_rate=0.01
)


@App.get('/')
def ping():
    return {"message": "pong"}

#: run: hypercorn app:App --reload
