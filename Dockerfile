FROM alpine:3.4

RUN addgroup router \
 && adduser -S -G router router

COPY ./src/github.com/Esjob-Cloud-DevOps/scheduler-router /go/src/github.com/Esjob-Cloud-DevOps/scheduler-router

RUN apk --update add ca-certificates \
 && apk --update add --virtual build-deps go git \
 && cd /go/src/github.com/Esjob-Cloud-DevOps/scheduler-router \
 && GOPATH=/go go build -o /bin/scheduler-router \
 && apk del --purge build-deps \
 && rm -rf /go/bin /go/pkg /var/cache/apk/*

USER router

EXPOSE 8899

ENTRYPOINT [ "/bin/scheduler-router" ]
