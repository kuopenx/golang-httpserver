FROM golang:1.18-alpine AS builder

ENV GO111MODULE=auto \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o httpserver .

FROM alpine
COPY --from=builder /build/httpserver .
EXPOSE 80
ENTRYPOINT ["./httpserver"]