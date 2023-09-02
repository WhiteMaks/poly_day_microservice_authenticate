FROM whitemaks/poly_day_migrator:release_0.0.1

ENV WORK_PATH=/usr/migrator

COPY ./migrations $WORK_PATH/migrations/

ENTRYPOINT ./entrypoint.sh