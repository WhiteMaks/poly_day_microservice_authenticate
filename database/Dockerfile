FROM whitemaks/poly_day_migrator:latest

ENV WORK_PATH=/usr/migrator

COPY ./migrations $WORK_PATH/migrations/

ENTRYPOINT ./entrypoint.sh