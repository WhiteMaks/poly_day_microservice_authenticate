FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add gcompat

EXPOSE 8080

ENV WORK_PATH=/usr/poly_day

COPY ./docker/entrypoint.sh $WORK_PATH/
COPY ./.bin/authenticate/poly_day_microservice_authenticate $WORK_PATH/

WORKDIR $WORK_PATH

RUN chmod 777 ./entrypoint.sh
RUN chmod 777 ./poly_day_microservice_authenticate

ENTRYPOINT ./entrypoint.sh