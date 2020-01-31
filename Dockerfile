FROM golang:alpine AS builder

RUN apk add --no-cache git openssh-client
ADD . /go/src/github.com/masterzen/winrm-cli
WORKDIR /go/src/github.com/masterzen/winrm-cli
RUN go mod tidy \
 && CGO_ENABLED="0" \
  GOOS="linux" \
  GOARCH="amd64" \
  go build -ldflags='-s -w' -o /go/bin/winrm

FROM scratch
COPY --from=builder /go/bin/winrm /winrm
ENTRYPOINT ["/winrm"]
