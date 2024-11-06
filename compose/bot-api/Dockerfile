FROM base-bot-py:latest

WORKDIR  /app

COPY ./deploy/compose/bot-api/start.sh ./cmd/start.sh
RUN chmod a+x cmd/start.sh


CMD ["sh", "./cmd/start.sh" ]
