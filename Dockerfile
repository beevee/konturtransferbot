FROM golang:1.10 as builder

WORKDIR /go/src/github.com/beevee/konturtransferbot

COPY . /go/src/github.com/beevee/konturtransferbot/

RUN go get github.com/kardianos/govendor
RUN govendor sync

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/konturtransferbot github.com/beevee/konturtransferbot/cmd/konturtransferbot


FROM alpine

WORKDIR /konturtransferbot

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY cmd/konturtransferbot/schedule.yml /konturtransferbot/
COPY --from=builder /go/src/github.com/beevee/konturtransferbot/build/konturtransferbot /konturtransferbot/

ENTRYPOINT ["./konturtransferbot"]
