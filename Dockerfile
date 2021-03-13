FROM golang:1.16.1-alpine AS build

WORKDIR $GOPATH/src/github.com/csachs/microwebdav

COPY . .

ARG CGO_ENABLED=0

RUN go build -o /microwebdav

FROM scratch

LABEL maintainer="sachs.christian@gmail.com"

COPY --from=build /microwebdav /microwebdav

ENV PATH= \
    MICROWEBDAV_LISTEN=:8000 \
    MICROWEBDAV_PATH=/data \
    MICROWEBDAV_USER=user

EXPOSE $MICROWEBDAV_LISTEN

VOLUME $MICROWEBDAV_PATH

ENTRYPOINT ["/microwebdav"]
