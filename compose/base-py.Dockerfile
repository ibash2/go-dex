FROM python:3.11-slim

COPY ./requirements.txt .

RUN pip install -r requirements.txt

COPY ../../src ./app
COPY ../../abis ./abis