FROM golang:1.18-alpine AS build

ENV GO111MODULE=auto \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o httpserver .

FROM scratch
COPY --from=build /build/httpserver /httpserver
EXPOSE 80
ENTRYPOINT ["/httpserver"]