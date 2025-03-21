FROM golang:1.23-alpine as golang

COPY go.mod go.sum /go/src/github.com/beevee/konturtransferbot/
WORKDIR /go/src/github.com/beevee/konturtransferbot

RUN go mod download
COPY . /go/src/github.com/beevee/konturtransferbot

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/konturtransferbot github.com/beevee/konturtransferbot/cmd/konturtransferbot

FROM alpine
RUN apk add --no-cache tzdata ca-certificates && update-ca-certificates
COPY cmd/konturtransferbot/schedule.yml /schedule.yml
COPY --from=golang /go/bin/konturtransferbot /konturtransferbot
ENTRYPOINT ["/konturtransferbot"]
