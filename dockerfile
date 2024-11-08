# builder image
FROM golang:1.22-alpine AS builder
RUN mkdir /build
WORKDIR /build
COPY go.mod go.sum ./
# you may use `GOPROXY` to speed it up in Mainland China.
# RUN  GOPROXY=https://goproxy.cn,direct go mod tidy
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wallet-server .

# final target image for multi-stage builds
FROM alpine:3.16
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/wallet-server .
COPY ./config/config.yml ./config.yml
ENTRYPOINT [ "./wallet-server" ]