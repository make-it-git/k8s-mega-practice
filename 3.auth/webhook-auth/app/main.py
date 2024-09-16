from fastapi import FastAPI, Request, HTTPException
from pydantic import BaseModel
from typing import Optional, List
import logging

logging.basicConfig(encoding='utf-8', level=logging.DEBUG)

logger = logging.getLogger(__name__)

app = FastAPI()

VALID_TOKEN = "secret"

class TokenReviewRequest(BaseModel):
    spec: dict
    apiVersion: str
    kind: str

class User(BaseModel):
    username: str
    groups: Optional[List[str]] = None  # Optional group memberships

class Status(BaseModel):
    authenticated: bool
    user: Optional[User] = None
    audiences: Optional[List[str]] = None  # Optional list of audiences

class TokenReviewResponse(BaseModel):
    apiVersion: str
    kind: str
    status: Status

@app.post("/validate-token")
async def validate_token(request: Request, body: TokenReviewRequest):
    logging.debug(body)
    token = body.spec['token']
    
    if token == VALID_TOKEN:
        response = TokenReviewResponse(
            apiVersion="authentication.k8s.io/v1",
            kind="TokenReview",
            status=Status(
                authenticated=True,
                user=User(
                    username="user2",
                    groups=["developers", "qa"],
                ),
                # audiences=["https://myserver.example.com"]
            )
        )
        logging.debug(response)
        return response

    response = TokenReviewResponse(
        apiVersion="authentication.k8s.io/v1",
        kind="TokenReview",
        status=Status(
            authenticated=False,
        )
    )
    logging.debug(response)
    return response

@app.get("/")
def read_root():
    return {"message": "Kubernetes Auth Webhook"}

