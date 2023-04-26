FROM alpine:latest

WORKDIR /app
COPY . /app

RUN apk add gcompat

EXPOSE 8081

ENTRYPOINT [ "/app/bin/auth" ]
