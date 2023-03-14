from pydantic import BaseModel

from typing import List

class Traffic(BaseModel):
    UserAgent: str
    Method: str
    Proto: str
    ProtoMajor: int
    ProtoMinor: int
    ContentLength: int
    TransferEncoding: List
    Close: bool
    RemoteAddr: str
    RequestURI: str
    Scheme: str
    Host: str
    Path: str
    IsMalicious: float = -1.0