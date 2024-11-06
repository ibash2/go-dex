#!/usr/bin/env bash
set -e
alembic upgrade head

DEFAULT_MODULE_NAME=application.entrypoints.api

MODULE_NAME=${MODULE_NAME:-$DEFAULT_MODULE_NAME}
VARIABLE_NAME=${VARIABLE_NAME:-app}
export APP_MODULE=${APP_MODULE:-"$MODULE_NAME:$VARIABLE_NAME"}

export WORKER_CLASS=${WORKER_CLASS:-"uvicorn.workers.UvicornWorker"}

gunicorn --forwarded-allow-ips "*" -k "$WORKER_CLASS" "$APP_MODULE" --bind 0.0.0.0:$APP_HTTP_PORT
