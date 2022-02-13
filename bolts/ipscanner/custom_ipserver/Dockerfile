FROM python:3.8-alpine

ENV IPSERVER_HOST 127.0.0.1
ENV IPSERVER_PORT 8000

RUN pip install fastapi uvicorn[standard]
RUN echo '# coding: utf-8\n \
import os\n \
import uvicorn\n \
from fastapi import FastAPI, Request\n \
from uuid import uuid4\n \
app = FastAPI()\n \
access_token = str(uuid4())\n \
@app.get("/ip/")\n \
async def get_ip(token: str, request: Request):\n \
    if token != access_token:\n \
        return {"status": -1, "message": "invalid token"}\n \
    return {"status": 0, "ip": request.client.host}\n \
host = os.getenv("IPSERVER_HOST", "127.0.0.1")\n \
port = os.getenv("IPSERVER_PORT", 8000)\n \
print("Access token is", access_token)\n \
uvicorn.run(app, host=host, port=int(port), proxy_headers=True, forwarded_allow_ips="*")\n \
' > /main.py

CMD python -u /main.py